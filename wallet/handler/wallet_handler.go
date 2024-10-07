package handler

import (
	pb "Assignment_Golang/user/proto/user_service/v1"
	"Assignment_Golang/wallet/entity"
	pbwallet "Assignment_Golang/wallet/proto/wallet_service/v1"
	"Assignment_Golang/wallet/service"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IWalletHandler interface {
	GetSaldoByUserID(c *gin.Context)
	GetHistoryTransactionByUserID(c *gin.Context)
	TopUpSaldoUser(c *gin.Context)
	DecreaseSaldoUser(c *gin.Context)
	TransferSaldoUser(c *gin.Context)
}

type WalletHandler struct {
	pbwallet.UnimplementedWalletServiceServer
	walletService service.IWalletService
	userClient    pb.UserServiceClient
}

func NewWalletHandler(walletService service.IWalletService, userClient pb.UserServiceClient) *WalletHandler {
	return &WalletHandler{walletService: walletService, userClient: userClient}
}

func (h *WalletHandler) GetHistoryTransactionByUserID(ctx context.Context, req *pbwallet.GetHistoryTransactionByUserIDRequest) (*pbwallet.GetHistoryTransactionByUserIDResponse, error) {
	_, err := h.getUserById(ctx, &pb.GetUserByIDRequest{Id: req.Id})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	historyTransactions, err := h.walletService.GetHistoryTransactionByUserID(ctx, int(req.Id))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var historyTransactionProto []*pbwallet.HistoryTransaction
	for _, historyTransaction := range historyTransactions {
		historyTransactionProto = append(historyTransactionProto, &pbwallet.HistoryTransaction{
			Id:              int32(historyTransaction.ID),
			IdUser:          int32(historyTransaction.IDUser),
			Name:            historyTransaction.Name,
			Saldo:           historyTransaction.Saldo,
			Amount:          historyTransaction.Amount,
			TransactionType: historyTransaction.TransactionType,
			TransactionDate: timestamppb.New(historyTransaction.TransactionDate),
		})
	}
	return &pbwallet.GetHistoryTransactionByUserIDResponse{
		HistoryTransactions: historyTransactionProto,
	}, nil
}

func (h *WalletHandler) GetSaldoByUserID(ctx context.Context, req *pbwallet.GetSaldoByUserIDRequest) (*pbwallet.GetSaldoByUserIDResponse, error) {
	_, err := h.getUserById(ctx, &pb.GetUserByIDRequest{Id: req.Id})
	if err != nil {
		wallet := &pbwallet.GetSaldoByUserIDResponse{
			SaldoWallet: &pbwallet.SaldoWallet{
				Id:     0,
				IdUser: 0,
				Name:   "",
				Saldo:  0,
			},
		}
		log.Println(err)
		return wallet, err
	}
	responsewallet, err := h.walletService.GetSaldoByUserID(ctx, int(req.Id))
	if err != nil {
		wallet := &pbwallet.GetSaldoByUserIDResponse{
			SaldoWallet: &pbwallet.SaldoWallet{
				Id:     0,
				IdUser: 0,
				Name:   "",
				Saldo:  0,
			},
		}
		log.Println("err")
		return wallet, err
	}

	res := &pbwallet.GetSaldoByUserIDResponse{
		SaldoWallet: &pbwallet.SaldoWallet{
			Id:     int32(responsewallet.ID),
			IdUser: int32(responsewallet.IDUser),
			Name:   responsewallet.Name,
			Saldo:  responsewallet.Saldo,
		},
	}
	return res, nil
}

func (h *WalletHandler) TopUpSaldoUser(ctx context.Context, req *pbwallet.TopUpSaldouserRequest) (*pbwallet.MutationResponse, error) {
	user, err := h.getUserById(ctx, &pb.GetUserByIDRequest{Id: req.IdUser})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	topupresponse, err := h.topUp(ctx, req, user)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbwallet.MutationResponse{
		Message: fmt.Sprintf("Success TopUp %f, current Saldo is %f", req.Amount, topupresponse.Saldo),
	}, nil
}

func (h *WalletHandler) DecreaseSaldoUser(ctx context.Context, req *pbwallet.DecreaseSaldouserRequest) (*pbwallet.MutationResponse, error) {
	_, err := h.getUserById(ctx, &pb.GetUserByIDRequest{Id: req.IdUser})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	decreaseresponse, err := h.decrease(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbwallet.MutationResponse{
		Message: fmt.Sprintf("Success debit Amount %f, current Saldo is %f", req.Amount, decreaseresponse.Saldo),
	}, nil
}

func (h *WalletHandler) TransferSaldoUser(ctx context.Context, req *pbwallet.TransferSaldouserRequest) (*pbwallet.MutationResponse, error) {
	_, err := h.getUserById(ctx, &pb.GetUserByIDRequest{Id: req.IdUserFrom})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	userto, err := h.getUserById(ctx, &pb.GetUserByIDRequest{Id: req.IdUserTo})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	updatedfromsaldo, err := h.decrease(ctx, &pbwallet.DecreaseSaldouserRequest{
		IdUser: req.IdUserFrom,
		Amount: req.Amount,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	updatedtosaldo, err := h.topUp(ctx, &pbwallet.TopUpSaldouserRequest{
		IdUser: req.IdUserTo,
		Amount: req.Amount,
	}, userto)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pbwallet.MutationResponse{
		Message: fmt.Sprintf("Success transfer amount %f from %s to %s", req.Amount, updatedfromsaldo.Name, updatedtosaldo.Name),
	}, nil
}

func (h *WalletHandler) getUserById(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := h.userClient.GetUserByID(ctx, req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if user.User.Id == 0 {
		log.Println("User Not Found")
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (h *WalletHandler) topUp(ctx context.Context, req *pbwallet.TopUpSaldouserRequest, user *pb.GetUserByIDResponse) (*entity.SaldoWallet, error) {
	lastHistoryTransaction, err := h.walletService.GetLastHistoryTransactionByUserID(ctx, int(req.IdUser))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if lastHistoryTransaction.ID == 0 {
		lastHistoryTransaction.IDUser = int(user.User.Id)
		lastHistoryTransaction.Name = user.User.Name
		lastHistoryTransaction.Saldo = 0
	}
	lastHistoryTransaction.Saldo = lastHistoryTransaction.Saldo + req.Amount
	lastHistoryTransaction.Amount = req.Amount
	topupresponse, err := h.walletService.TopUpSaldoUser(ctx, lastHistoryTransaction)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return topupresponse, nil
}

func (h *WalletHandler) decrease(ctx context.Context, req *pbwallet.DecreaseSaldouserRequest) (*entity.SaldoWallet, error) {
	lastHistoryTransaction, err := h.walletService.GetLastHistoryTransactionByUserID(ctx, int(req.IdUser))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if lastHistoryTransaction.ID == 0 {
		log.Println("User Doesn't have Saldo")
		return nil, fmt.Errorf("failed to debit wallet, user doesn't have saldo")
	}
	lastHistoryTransaction.Amount = req.Amount
	decreaseresponse, err := h.walletService.DecreaseSaldoUser(ctx, lastHistoryTransaction)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return decreaseresponse, nil
}
