package models

import (
	. "app/common"
	"app/services"
	"sync"
)

type UserBalance struct {
	Id      int        `json:"id"`
	Name    string     `json:"name"`
	Balance int        `json:"balance"`
	mu      sync.Mutex `json:"-"`
}

func (m *UserBalance) GetById(id int) *UserBalance {
	error := services.DB.Where("id = ?", id).First(m).Error
	if error != nil {
		Log.Printf("Unable to get user balance!  %s", error, id)
	}

	return m
}

func (m *UserBalance) AddAmount(summ int) (*UserBalance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	amount := m.Balance + summ
	error := services.DB.Model(m).Update("balance", amount).Error
	if error != nil {
		Log.Printf("Unable to update user balance!  %s", error, summ, m, amount)
	}
	updated := m.GetById(m.Id)

	return updated, error
}

func (m *UserBalance) WithdrawAmount(summ int) (*UserBalance, error) {
	var amount int
	var error error
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Balance < summ {
		amount = m.Balance
	} else {
		amount = m.Balance - summ
	}
	error = services.DB.Model(m).Update("balance", amount).Error
	if error != nil {
		Log.Printf("Unable to withdraw user balance!  %s", error, summ, m)
	}
	updated := m.GetById(m.Id)

	return updated, error
}

func CreateUserBalance(id, balance int, name string) *UserBalance{
	model := UserBalance{
		Id:id,
		Balance:balance,
		Name:name,
	}
	services.DB.Create(&model)

	return  &model
}

func (m UserBalance) TableName() string {
	return "user_balance"
}
