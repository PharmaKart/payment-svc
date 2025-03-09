package repositories

import (
	"fmt"

	"github.com/PharmaKart/payment-svc/internal/models"
	"github.com/PharmaKart/payment-svc/pkg/errors"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	StorePayment(payment *models.Payment) error
	GetPaymentByOrderID(orderID string) (*models.Payment, error)
	GetPaymentByTransactionID(transactionID string) (*models.Payment, error)
	GetPayment(paymentID string) (*models.Payment, error)
	UpdatePaymentStatus(orderID string, status string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) StorePayment(payment *models.Payment) error {
	if err := r.db.Create(payment).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (r *paymentRepository) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("Payment for order ID '%s' not found", orderID))
		}
		return nil, errors.NewInternalError(err)
	}
	return &payment, nil
}

func (r *paymentRepository) GetPaymentByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("Payment with transaction ID '%s' not found", transactionID))
		}
		return nil, errors.NewInternalError(err)
	}
	return &payment, nil
}

func (r *paymentRepository) GetPayment(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("id = ?", paymentID).First(&payment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError(fmt.Sprintf("Payment with ID '%s' not found", paymentID))
		}
		return nil, errors.NewInternalError(err)
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentStatus(orderID string, status string) error {
	result := r.db.Model(&models.Payment{}).Where("order_id = ?", orderID).Update("status", status)

	if result.Error != nil {
		return errors.NewInternalError(result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("Payment for order ID '%s' not found", orderID))
	}

	return nil
}
