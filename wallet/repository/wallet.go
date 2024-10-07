package repository

import (
	walletEntity "Assignment_Golang/wallet/entity"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

type GormDBFace interface {
	WithContext(ctx context.Context) *gorm.DB
	Create(value interface{}) *gorm.DB
	First(value interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
}

type IWalletRepository interface {
	GetSaldoByUserID(ctx context.Context, id int) (walletEntity.SaldoWallet, error)
	GetHistoryTransactionByUserID(ctx context.Context, id int) ([]walletEntity.HistoryTransaction, error)
	GetLastHistoryTransactionByUserID(ctx context.Context, id int) (walletEntity.HistoryTransaction, error)
	TopUpSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error)
	DecreaseSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error)
	// TransferSaldoUser(ctx context.Context, updatedwalletfrom walletEntity.HistoryTransaction, updatedwalletto walletEntity.HistoryTransaction) error
}

// userRepository adalah implementasi dari IUserRepository yang menggunakan slice untuk menyimpan data user
type walletRepository struct {
	db GormDBFace
}

// NewRepository membuat instance baru dari userRepository
func NewWalletRepository(db GormDBFace) IWalletRepository {
	return &walletRepository{
		db: db,
	}
}

func (r *walletRepository) GetSaldoByUserID(ctx context.Context, id int) (walletEntity.SaldoWallet, error) {
	var saldowallet walletEntity.SaldoWallet
	query := "SELECT id, id_user, name, saldo FROM wallet WHERE id_user = $1 ORDER BY transaction_date DESC LIMIT 1"
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&saldowallet).Error; err != nil {
		log.Printf("Error Get Saldo: %v\n", err)
		return walletEntity.SaldoWallet{}, err
	}
	if saldowallet.ID == 0 {
		return walletEntity.SaldoWallet{}, gorm.ErrRecordNotFound
	}
	return saldowallet, nil
}

// GetHistoryTransactionByUserID mengembalikan semua transaksi yang dilakukan pengguna
func (r *walletRepository) GetHistoryTransactionByUserID(ctx context.Context, id int) ([]walletEntity.HistoryTransaction, error) {
	var historyTransaction []walletEntity.HistoryTransaction
	query := "SELECT id, id_user, name, saldo, amount, transaction_type, transaction_date FROM wallet WHERE id_user = $1"
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&historyTransaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return historyTransaction, nil
		}
		log.Printf("Error getting history transaction: %v\n", err)
		return historyTransaction, err
	}
	return historyTransaction, nil
}

func (r *walletRepository) GetLastHistoryTransactionByUserID(ctx context.Context, id int) (walletEntity.HistoryTransaction, error) {
	var historyTransaction walletEntity.HistoryTransaction
	query := "SELECT id, id_user, name, saldo, amount, transaction_type, transaction_date FROM wallet WHERE id_user = $1 ORDER BY transaction_date DESC LIMIT 1"
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&historyTransaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return walletEntity.HistoryTransaction{ID: 0}, nil
		}
		log.Printf("Error getting history transaction: %v\n", err)
		return walletEntity.HistoryTransaction{ID: 0}, err
	}
	return historyTransaction, nil
}

func (r *walletRepository) TopUpSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error) {
	var createdId int
	var updatedsaldo walletEntity.SaldoWallet
	query := "INSERT INTO wallet(id_user, name, saldo, amount, transaction_type, transaction_date) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id"
	if err := r.db.WithContext(ctx).Raw(query, updatedwallet.IDUser, updatedwallet.Name, updatedwallet.Saldo, updatedwallet.Amount, "TopUp").Scan(&createdId).Error; err != nil {
		log.Printf("Error TopUp saldo Wallet: %v\n", err)
		return nil, err
	}
	updatedsaldo.ID = createdId
	updatedsaldo.IDUser = updatedwallet.IDUser
	updatedsaldo.Name = updatedwallet.Name
	updatedsaldo.Saldo = updatedwallet.Saldo
	return &updatedsaldo, nil
}

func (r *walletRepository) DecreaseSaldoUser(ctx context.Context, updatedwallet walletEntity.HistoryTransaction) (*walletEntity.SaldoWallet, error) {
	var createdId int
	var updatedsaldo walletEntity.SaldoWallet
	query := "INSERT INTO wallet(id_user, name, saldo, amount, transaction_type, transaction_date) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id"
	if err := r.db.WithContext(ctx).Raw(query, updatedwallet.IDUser, updatedwallet.Name, updatedwallet.Saldo, updatedwallet.Amount, "Decrease").Scan(&createdId).Error; err != nil {
		log.Printf("Error Debit saldo Wallet: %v\n", err)
		return nil, err
	}
	updatedsaldo.ID = createdId
	updatedsaldo.IDUser = updatedwallet.IDUser
	updatedsaldo.Name = updatedwallet.Name
	updatedsaldo.Saldo = updatedwallet.Saldo
	return &updatedsaldo, nil
}

// func (r *walletRepository) TransferSaldoUser(ctx context.Context, updatedwalletfrom walletEntity.HistoryTransaction, updatedwalletto walletEntity.HistoryTransaction) error {
// 	var createdId int
// 	query := "INSERT INTO wallet(id_user, name, saldo, amount, transaction_type, transaction_date) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id"
// 	if err := r.db.WithContext(ctx).Raw(query, updatedwalletfrom.IDUser, updatedwalletfrom.Name, updatedwalletfrom.Saldo, updatedwalletfrom.Amount, "Decrease").Scan(&createdId).Error; err != nil {
// 		log.Printf("Error Debit saldo Wallet: %v\n", err)
// 		return err
// 	}
// 	if err := r.db.WithContext(ctx).Raw(query, updatedwalletto.IDUser, updatedwalletto.Name, updatedwalletto.Saldo, updatedwalletto.Amount, "TopUp").Scan(&createdId).Error; err != nil {
// 		log.Printf("Error TopUp saldo Wallet: %v\n", err)
// 		return err
// 	}
// 	return nil
// }
