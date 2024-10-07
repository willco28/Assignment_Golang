package service

import (
	walletEntity "Assignment_Golang/wallet/entity"
	"Assignment_Golang/wallet/repository"
	"context"
	"fmt"
)

// service ==> buat interface --> struct --> constructor nya --> bikin function implement method dari interface nya
type IWalletService interface { //interface dengan method yang berkaitan dengan wallet
	GetSaldoByUserID(ctx context.Context, id int) (walletEntity.SaldoWallet, error)
	GetHistoryTransactionByUserID(ctx context.Context, id int) ([]walletEntity.HistoryTransaction, error)
	GetLastHistoryTransactionByUserID(ctx context.Context, id int) (walletEntity.HistoryTransaction, error)
	TopUpSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error)
	DecreaseSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error)
	// TransferSaldoUser(ctx context.Context, updatedwalletfrom walletEntity.HistoryTransaction, updatedwalletto walletEntity.HistoryTransaction) error
}

type walletService struct {
	walletRepo repository.IWalletRepository //bikin struct dengan balikan interface repo
}

func NewWalletService(walletRepo repository.IWalletRepository) IWalletService {
	//bikin constructor dengan balikan wallet service yang isinya wallet repo
	return &walletService{walletRepo: walletRepo}
}

func (s *walletService) GetSaldoByUserID(ctx context.Context, id int) (walletEntity.SaldoWallet, error) {
	//kalau ada validasi ditambahkan disini
	saldoUser, err := s.walletRepo.GetSaldoByUserID(ctx, id)
	if err != nil {
		return walletEntity.SaldoWallet{}, fmt.Errorf("error get saldo user: %v", err)
	}
	return saldoUser, nil
}

func (s *walletService) GetHistoryTransactionByUserID(ctx context.Context, id int) ([]walletEntity.HistoryTransaction, error) {
	historyTransaction, err := s.walletRepo.GetHistoryTransactionByUserID(ctx, id)
	if err != nil {
		return []walletEntity.HistoryTransaction{}, fmt.Errorf("failed to retreive transaction history: %v", err)
	}
	return historyTransaction, nil
}

func (s *walletService) GetLastHistoryTransactionByUserID(ctx context.Context, id int) (walletEntity.HistoryTransaction, error) {
	historyTransaction, err := s.walletRepo.GetLastHistoryTransactionByUserID(ctx, id)
	if err != nil {
		return walletEntity.HistoryTransaction{ID: 0}, fmt.Errorf("failed to retreive transaction history: %v", err)
	}
	return historyTransaction, nil
}

func (s *walletService) TopUpSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error) {
	//kalau ada validasi ditambahkan disini
	newsaldo, err := s.walletRepo.TopUpSaldoUser(ctx, updatedwallet)
	if err != nil {
		return nil, fmt.Errorf("failed to TopUp saldo Wallet: %v", err)
	}
	return newsaldo, nil
}

func (s *walletService) DecreaseSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error) {
	//kalau ada validasi ditambahkan disini
	currentsaldo, err := s.walletRepo.GetSaldoByUserID(ctx, updatedwallet.IDUser)
	if err != nil {
		return nil, fmt.Errorf("error get saldo user: %v", err)
	}
	if currentsaldo.Saldo < updatedwallet.Amount {
		return nil, fmt.Errorf("failed to Debit saldo Wallet: Not Enough Saldo")
	}
	updatedwallet.Saldo = currentsaldo.Saldo - updatedwallet.Amount
	newsaldo, err := s.walletRepo.DecreaseSaldoUser(ctx, updatedwallet)
	if err != nil {
		return nil, fmt.Errorf("failed to Debit saldo Wallet: %v", err)
	}
	return newsaldo, nil
}

// func (s *walletService) TransferSaldoUser(ctx context.Context, updatedwalletfrom walletEntity.HistoryTransaction, updatedwalletto walletEntity.HistoryTransaction) error {
// 	//bikin function untuk mengembalikan wallet repo
// 	currentsaldo, err := s.walletRepo.GetSaldoByUserID(ctx, updatedwalletfrom.IDUser)
// 	if err != nil {
// 		return fmt.Errorf("failed to retreive all users: %v", err)
// 	}
// 	if currentsaldo.Saldo < updatedwalletfrom.Amount {
// 		return fmt.Errorf("failed to Debit saldo Wallet: Not Enough Saldo")
// 	}
// 	updatedwalletfrom.Saldo = currentsaldo.Saldo - updatedwalletfrom.Amount
// 	err = s.walletRepo.TransferSaldoUser(ctx, updatedwalletfrom, updatedwalletto)
// 	if err != nil {
// 		return fmt.Errorf("failed to Transfer saldo Wallet: %v", err)
// 	}
// 	return nil
// }
