package room

import (
	"encoding/json"
	"errors"
	"lock-stock-v2/handlers/http/helpers"
	gameRepository "lock-stock-v2/internal/domain/game/repository"
	gameService "lock-stock-v2/internal/domain/game/service"
	roomModel "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room/repository"
	"lock-stock-v2/internal/domain/room/service"
	"lock-stock-v2/internal/domain/room_user/service"
	userRepository "lock-stock-v2/internal/domain/user/repository"
	"lock-stock-v2/internal/websocket"
	"log"
	"net/http"
)

type RoomHandler struct {
	joinRoomService          *services.JoinRoomService
	roomUserService          *services.RoomUserService
	startGameService         *service.StartGameService
	sendAnswerService        *gameService.SendAnswer
	roomRepository           repository.RoomRepository
	userRepository           userRepository.UserRepository
	createBet                *gameService.CreateBetService
	playerRepository         gameRepository.PlayerRepository
	roundRepository          gameRepository.RoundRepository
	betRepository            gameRepository.BetRepository
	gameRepository           gameRepository.GameRepository
	roundPlayerLogRepository gameRepository.RoundPlayerLogRepository
	webSocket                websocket.Manager
}

func NewRoomHandler(
	u *services.JoinRoomService,
	roomRepository repository.RoomRepository,
	userRepository userRepository.UserRepository,
	roomUserService *services.RoomUserService,
	startGameService *service.StartGameService,
	createBet *gameService.CreateBetService,
	playerRepository gameRepository.PlayerRepository,
	roundRepository gameRepository.RoundRepository,
	betRepository gameRepository.BetRepository,
	gameRepository gameRepository.GameRepository,
	roundPlayerLogRepository gameRepository.RoundPlayerLogRepository,
	webSocket websocket.Manager,
	sendAnswerService *gameService.SendAnswer,
) *RoomHandler {
	return &RoomHandler{
		joinRoomService:          u,
		roomRepository:           roomRepository,
		roomUserService:          roomUserService,
		startGameService:         startGameService,
		userRepository:           userRepository,
		createBet:                createBet,
		playerRepository:         playerRepository,
		roundRepository:          roundRepository,
		betRepository:            betRepository,
		gameRepository:           gameRepository,
		roundPlayerLogRepository: roundPlayerLogRepository,
		webSocket:                webSocket,
		sendAnswerService:        sendAnswerService,
	}
}

func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.roomRepository.GetPending()
	if err != nil {
		respondWithError(w, "Failed to get rooms", err, http.StatusInternalServerError)
		return
	}
	responseData := make([]RoomResponse, 0)

	for _, room := range rooms {
		var responseRoomUid = room.Uid()
		responseData = append(responseData, RoomResponse{
			RoomUid: &responseRoomUid,
		})
	}

	respondWithJSON(w, http.StatusOK, responseData)
}

func (h *RoomHandler) StartGame(w http.ResponseWriter, r *http.Request, roomId string) {
	room, err := helpers.GetRoomById(h.roomRepository, roomId)
	if room == nil || err != nil {
		respondWithError(w, "Error getting room", nil, http.StatusBadRequest)
		return
	}
	if room.Status() == roomModel.StatusStarted {
		respondWithError(w, "Room already started", nil, http.StatusBadRequest)
		return
	}

	user, err := helpers.GetUserFromRequest(r)
	if err != nil {
		var userErr *helpers.UserNotFoundError
		ok := errors.As(err, &userErr)
		if ok {
			respondWithError(w, err.Error(), nil, userErr.Code)
		} else {
			respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
		}
		return
	}

	roomUsers, err := h.roomUserService.GetUsersByRoom(room)
	if err != nil {
		http.Error(w, "Failed to get users in room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	isUserInRoom := h.roomUserService.IsUserInRoom(roomUsers, user)

	if !isUserInRoom {
		http.Error(w, "RoomUser is not in the room", http.StatusForbidden)
		return
	}

	req := service.StartGameRequest{
		Room: room,
		User: user,
	}

	if err := h.startGameService.StartGame(req); err != nil {
		respondWithError(w, "Failed to start game", err, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Game started"})
}

func (h *RoomHandler) SendAnswer(w http.ResponseWriter, r *http.Request, params SendAnswerParams) {
	user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
	if err != nil {
		var userErr *helpers.UserNotFoundError
		ok := errors.As(err, &userErr)
		if ok {
			respondWithError(w, err.Error(), nil, userErr.Code)
		} else {
			respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
		}
		return
	}
	if user == nil {
		log.Println("User not found")
		return
	}
	game, err := h.gameRepository.FindByUser(user)
	if err != nil {
		log.Printf("Game by user %s not found", user.Uid())
		respondWithError(w, "Error sending answer", err, http.StatusInternalServerError)
		return
	}
	round, err := h.roundRepository.FindLastByGame(game)
	if err != nil {
		log.Printf("Round by game %s not found", game.Uid())
		return
	}
	var nwkRawAnswer NwkRawAnswer
	if err = json.NewDecoder(r.Body).Decode(&nwkRawAnswer); err != nil {
		respondWithError(w, "invalid decoding request body", err, http.StatusBadRequest)
		return
	}
	roundPlayerLog, err := h.roundPlayerLogRepository.FindByRoundAndUser(round, user)
	if err != nil {
		respondWithError(w, "round player log not found", err, http.StatusBadRequest)
		return
	}
	if nil != roundPlayerLog {
		answer := uint(nwkRawAnswer.Value)
		err = h.sendAnswerService.SendAnswer(roundPlayerLog, answer)
		if err != nil {
			respondWithError(w, "Failed on save answer", err, http.StatusBadRequest)
			return
		}
	}
	respondWithJSON(w, http.StatusOK, "SUCCESS")
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request, roomId string, params JoinRoomParams) {
	room, err := helpers.GetRoomById(h.roomRepository, roomId)
	if err != nil {
		var roomErr *helpers.RoomNotFoundError
		ok := errors.As(err, &roomErr)
		if ok {
			respondWithError(w, err.Error(), nil, roomErr.Code)
		} else {
			respondWithError(w, "Error getting room", err, http.StatusInternalServerError)
		}
		return
	}

	user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
	if err != nil {
		var userErr *helpers.UserNotFoundError
		ok := errors.As(err, &userErr)
		if ok {
			respondWithError(w, err.Error(), nil, userErr.Code)
		} else {
			respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
		}
		return
	}

	domainReq := services.JoinRoomRequest{User: user, Room: room}
	if err := h.joinRoomService.JoinRoom(domainReq); err != nil {
		respondWithError(w, "Failed to join room", err, http.StatusInternalServerError)
		return
	}

	roomUsers, err := h.roomUserService.GetUsersByRoom(room)
	if err != nil {
		respondWithError(w, "Failed to get users in room", err, http.StatusInternalServerError)
		return
	}

	response := make([]JoinRoomResponse, 0, len(roomUsers))
	for _, ru := range roomUsers {
		response = append(response, JoinRoomResponse{
			RoomId:   ru.Room().Uid(),
			UserId:   ru.User().Uid(),
			UserName: ru.User().Name(),
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *RoomHandler) MakeBet(w http.ResponseWriter, r *http.Request, params MakeBetParams) {
	user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
	if err != nil {
		var userErr *helpers.UserNotFoundError
		ok := errors.As(err, &userErr)
		if ok {
			respondWithError(w, err.Error(), nil, userErr.Code)
		} else {
			respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
		}
		return
	}
	var nwkRawBet NwkRawBet
	if err := json.NewDecoder(r.Body).Decode(&nwkRawBet); err != nil {
		respondWithError(w, "invalid request body", err, http.StatusBadRequest)
		return
	}
	room, err := h.roomRepository.FindById(nwkRawBet.RoomId)
	if err != nil {
		log.Printf("error finding room by ID %s: %v", nwkRawBet.RoomId, err)
		return
	}

	player, err := h.playerRepository.FindByUserAndRoom(user, room)
	if err != nil {
		log.Printf("error finding player for user %s in room %s: %v", user.Uid(), room.Uid(), err)
		return
	}

	round, err := h.roundRepository.FindLastByGame(player.Game())
	if err != nil {
		log.Printf("error finding last round for game %s: %v", player.Game().Uid(), err)
		return
	}

	_, err = h.createBet.CreateBet(player, nwkRawBet.Amount, round)
	if err != nil {
		log.Println("error creating bet:", err)
		respondWithError(w, "error creating bet", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, "success")
}

func respondWithError(w http.ResponseWriter, message string, err error, statusCode int) {
	errorMessage := message
	if err != nil {
		errorMessage += ": " + err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
	if err != nil {
		log.Println("Failed to encode error response:", err)
		return
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
