package room

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"lock-stock-v2/external/transaction"
	"lock-stock-v2/external/websocket"
	"lock-stock-v2/handlers/http/helpers"
	gameRepository "lock-stock-v2/internal/domain/game/repository"
	gameService "lock-stock-v2/internal/domain/game/service"
	roomModel "lock-stock-v2/internal/domain/room/model"
	"lock-stock-v2/internal/domain/room/repository"
	"lock-stock-v2/internal/domain/room/service"
	"lock-stock-v2/internal/domain/room_user/service"
	userRepository "lock-stock-v2/internal/domain/user/repository"
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
	transactionManager       transaction.TransactionManager
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
	transactionManager transaction.TransactionManager,
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
		transactionManager:       transactionManager,
	}
}

var RoomAlreadyStarted = errors.New("room already started")
var RoomNotFound = errors.New("room not found")

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
	ctx := r.Context()

	err := h.transactionManager.Run(ctx, func(tx pgx.Tx) error {
		room, err := helpers.GetRoomById(h.roomRepository, roomId)
		if nil == room {
			log.Printf("Room by id %s not found", roomId)
			return RoomNotFound
		}
		if err != nil {
			log.Printf("error getting room: %s", err.Error())
			return err
		}
		if room.Status() == roomModel.StatusStarted {
			log.Println("room already started")
			return RoomAlreadyStarted
		}

		user, err := helpers.GetUserFromRequest(r)
		if err != nil {
			return err
		}

		req := service.StartGameRequest{
			Room: room,
			User: user,
		}

		return h.startGameService.StartGame(ctx, tx, req)
	})

	if err != nil {
		respondWithError(w, "Failed to start game", err, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Game started"})
}

func (h *RoomHandler) SendAnswer(w http.ResponseWriter, r *http.Request, params SendAnswerParams) {
	ctx := r.Context()

	err := h.transactionManager.Run(ctx, func(tx pgx.Tx) error {
		user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
		if err != nil || nil == user {
			log.Println("Error getting user")
			return err
		}
		game, err := h.gameRepository.FindByUser(user)
		if err != nil || game == nil {
			log.Printf("Game by user %s not found", user.Uid())
			return err
		}
		round, err := h.roundRepository.FindLastByGame(game)
		if err != nil {
			log.Printf("Round by game %s not found", game.Uid())
			return err
		}
		var nwkRawAnswer NwkRawAnswer
		if err = json.NewDecoder(r.Body).Decode(&nwkRawAnswer); err != nil {
			log.Println("error decoding answer")
			return err
		}
		roundPlayerLog, err := h.roundPlayerLogRepository.FindByRoundAndUser(round, user)
		if err != nil {
			return err
		}
		if nil != roundPlayerLog {
			answer := uint(nwkRawAnswer.Value)
			err = h.sendAnswerService.SendAnswer(ctx, tx, roundPlayerLog, answer)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if nil != err {
		respondWithError(w, "Failed on send answer", err, http.StatusBadRequest)
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
	ctx := r.Context()
	err := h.transactionManager.Run(ctx, func(tx pgx.Tx) error {
		user, err := helpers.GetUserFromString(params.Authorization, h.userRepository)
		if err != nil {
			var userErr *helpers.UserNotFoundError
			ok := errors.As(err, &userErr)
			if ok {
				respondWithError(w, err.Error(), nil, userErr.Code)
			} else {
				respondWithError(w, "Error getting user", err, http.StatusInternalServerError)
			}
			return err
		}
		var nwkRawBet NwkRawBet
		if err := json.NewDecoder(r.Body).Decode(&nwkRawBet); err != nil {
			respondWithError(w, "invalid request body", err, http.StatusBadRequest)
			return err
		}
		room, err := h.roomRepository.FindById(nwkRawBet.RoomId)
		if err != nil {
			log.Printf("error finding room by ID %s: %v", nwkRawBet.RoomId, err)
			return err
		}

		player, err := h.playerRepository.FindByUserAndRoom(user, room)
		if err != nil {
			log.Printf("error finding player for user %s in room %s: %v", user.Uid(), room.Uid(), err)
			return err
		}

		round, err := h.roundRepository.FindLastByGame(player.Game())
		if err != nil {
			log.Printf("error finding last round for game %s: %v", player.Game().Uid(), err)
			return err
		}

		bet, err := h.createBet.CreateBet(ctx, tx, player, nwkRawBet.Amount, round)
		if nil == bet {
			log.Println("Error creating bet")
			return err
		}
		return nil
	})

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
