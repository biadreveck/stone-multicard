package models

import (
	"errors"	
	"time"
)

type Wallet struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
    UpdatedAt time.Time `json:"-"`
    // DeletedAt *time.Time `json:"-"`

	Name			string `json:"name" gorm:"not null"`
	UserLimit		float64 `json:"userLimit"`
	MaxLimit		float64 `json:"maxLimit" gorm:"-"`
	AvailableCredit	float64 `json:"availableCredit" gorm:"-"`

	UserId	uint64 `json:"-"`
	Cards	[]Card `json:"cards"`
}

func (w *Wallet) Create() error {
	if err := w.validateCreation(); err != nil {
		return err
	}
	return database.Create(&w).Error
}

func (w *Wallet) Update() error {
	if err := w.validateUpdate(); err != nil {
		return err
	}
	updateFields := Wallet{Name: w.Name, UserLimit: w.UserLimit}
	return database.Set("gorm:save_associations", false).Model(&w).Updates(updateFields).Error
}

func (w Wallet) FechtAllFromUser() []Wallet {
	var wallets []Wallet
	database.Where("user_id = ?", w.UserId).Preload("Cards").Find(&wallets)
	for i, _ := range wallets {
		wallets[i].fillNonDBFields()
    }
	return wallets
}

func (w *Wallet) Get() bool {
	query := database.Preload("Cards").First(&w, w.ID)
	queryOk := query.Error == nil && !query.RecordNotFound()
	if queryOk {
		w.fillNonDBFields()
	}
	return queryOk
}

func (w *Wallet) Delete() error {
	return database.Delete(&w).Error
}

func (w *Wallet) Purchase(value float64) error {
	cards := w.Cards
	orderByDueDate := func(c1, c2 *Card) bool {
		now := time.Now()
		return c1.DueDayJson.NextDueDate().Sub(now) > c2.DueDayJson.NextDueDate().Sub(now)
	}
	orderByLimit := func(c1, c2 *Card) bool {
		return c1.AvailableCredit > c2.AvailableCredit
	}
	orderedBy(orderByDueDate, orderByLimit).Sort(cards)

	tx := database.Begin()
	for _, card := range cards {
		breakLoop := false;
		if card.AvailableCredit >= value {
			card.AvailableCredit -= value
			breakLoop = true;
		} else {
			card.AvailableCredit = 0
			value -= card.AvailableCredit
		}

		err := tx.Model(&card).UpdateColumn("available_credit", card.AvailableCredit).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if breakLoop {
			break
		}
    }
	tx.Commit()

	return nil
}

func (w *Wallet) validateCreation() error {
	if w.Name == "" {
		return errors.New("The field 'name' from Wallet must have a value")
	}
	maxLimit := 0.0
	if len(w.Cards) > 0 {
		for i, _ := range w.Cards {
	        if err := w.Cards[i].validateCreation(); err != nil {
	        	return err
	        }			
	        maxLimit += w.Cards[i].Limit
	    }
	}
	if w.UserLimit != 0 {
		if w.UserLimit > maxLimit {
			w.UserLimit = maxLimit
		} else if w.UserLimit < 0 {
			w.UserLimit = 0
		}
	}
	return nil
}
func (w *Wallet) validateUpdate() error {
	_wallet := &Wallet{ ID: w.ID }
	if !_wallet.Get() {
		return errors.New("Wallet not found")
	}
	maxLimit := _wallet.MaxLimit
	if len(w.Cards) > 0 {
		for _, c := range w.Cards {
	        if err := c.validateCreation(); err != nil {
	        	return err
	        }
	        maxLimit += c.Limit
	    }
	}
	if w.UserLimit != 0 {
		if w.UserLimit > maxLimit {
			w.UserLimit = maxLimit
		} else if w.UserLimit < 0 {
			w.UserLimit = 0
		}
	}
	return nil
}
func (w *Wallet) fillNonDBFields() {
	var usedCredit float64
	w.MaxLimit = 0
	for i, _ := range w.Cards {
		w.MaxLimit += w.Cards[i].Limit
		usedCredit += (w.Cards[i].Limit - w.Cards[i].AvailableCredit)		
		w.Cards[i].fillNonDBFields()
    }
	w.AvailableCredit = w.UserLimit - usedCredit
}