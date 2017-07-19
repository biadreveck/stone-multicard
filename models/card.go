package models

import (
	"errors"
	"sort"
	"time"

	"github.com/jinzhu/gorm"
)

type Card struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
    UpdatedAt time.Time `json:"-"`
    // DeletedAt *time.Time `json:"-"`

	Number			string `json:"number" gorm:"not null"`
	DueDay			*time.Time `json:"-" gorm:"not null"`
	DueDayJson		*dueDay `json:"dueDay" gorm:"-"`
	ExpirationDate	*time.Time `json:"-" gorm:"not null"`
	ExpirationDateJson	*expirationDate `json:"expirationDate" gorm:"-"`
	CVV				string `json:"cvv" gorm:"not null"`
	Limit			float64 `json:"limit" gorm:"not null"`
	AvailableCredit	float64 `json:"availableCredit"`

	// DueDay		*dueDay `json:"dueDay" gorm:"not null"`
	// ExpirationDate	*expirationDate `json:"expirationDate" gorm:"not null"`

	WalletId	uint64 `json:"-"`
}

func (c *Card) Create() error {	
	if err := c.validateCreation(); err != nil {
		return err
	}
	return database.Create(&c).Error
}

func (c *Card) Update() error {	
	if err := c.validateUpdate(); err != nil {
		return err
	}
	updateFields := Card{DueDay: c.DueDay, Limit: c.Limit, AvailableCredit: c.AvailableCredit}
	return database.Model(&c).Updates(updateFields).Error
}

func (c Card) FechtAll() []Card {
	var cards []Card
	database.Find(&cards)
	for i, _ := range cards {		
		cards[i].fillNonDBFields()
    }
	return cards
}

func (c Card) FechtAllFromWallet() []Card {
	var cards []Card
	database.Where("wallet_id = ?", c.WalletId).Find(&cards)
	for i, _ := range cards {		
		cards[i].fillNonDBFields()
    }
	return cards
}

func (c *Card) Get() bool {
	query := database.First(&c, c.ID)
	queryOk := query.Error == nil && !query.RecordNotFound()
	if queryOk {
		c.fillNonDBFields()
	}	
	return queryOk
}

func (c *Card) Delete() error {
	if !c.Get() {
		return nil
	}

	tx := database.Begin()

	if err := tx.Delete(&c).Error; err != nil {
		tx.Rollback()
		return err
	}

	wallet := &Wallet{ ID: c.WalletId }
	walletQuery := tx.Preload("Cards").First(&wallet, wallet.ID)
	if walletQuery.Error != nil || walletQuery.RecordNotFound() {
		return nil
	}
	wallet.fillNonDBFields()

	if wallet.UserLimit > wallet.MaxLimit {
		updateFields := Wallet{UserLimit: wallet.MaxLimit}
		if err := tx.Set("gorm:save_associations", false).Model(&wallet).Updates(updateFields).Error; err != nil {
			tx.Rollback()
			return err
		}		
	}

	tx.Commit()

	return nil
}

func (c *Card) ResetCredit() error {
	// dueDateYear, dueDateMonth, dueDateDay := c.DueDay.Date()
	// if dueDateMonth == time.December {
	// 	dueDateMonth = time.January
	// 	dueDateYear += 1
	// } else {
	// 	dueDateMonth += 1
	// }
	// newDueDate := time.Date(dueDateYear, dueDateMonth, dueDateDay, 0, 0, 0, 0, time.UTC)

	// updateFields := map[string]interface{}{"due_date": newDueDate, "available_credit": gorm.Expr("limit")}
	// return database.Model(&c).Updates(updateFields).Error
	
	// TODO: mudar para atualização de um campo só no gorm
	updateFields := map[string]interface{}{"available_credit": gorm.Expr("limit")}
	return database.Model(&c).Updates(updateFields).Error
}

func (c *Card) validateCreation() error {
	if c.Number == "" {
		return errors.New("The field 'number' from Card must have a value")
	}	
	if c.DueDayJson == nil {
		return errors.New("The field 'dueDay' from Card must have a value")
	}	
	if c.ExpirationDateJson == nil {
		return errors.New("The field 'expirationDate' from Card must have a value")
	}
	if c.CVV == "" {
		return errors.New("The field 'cvv' from Card must have a value")
	}
	if c.Limit == 0 {
		return errors.New("The field 'limit' from Card must have a value")
	}
	if c.AvailableCredit != 0 {
		if c.AvailableCredit > c.Limit {
			c.AvailableCredit = c.Limit
		} else if c.AvailableCredit < 0 {
			c.AvailableCredit = 0
		}
	}

	if len(c.Number) != 16 {
		return errors.New("The field 'number' from Card is invalid, must have 16 digits")
	}
	if len(c.CVV) != 3 {
		return errors.New("The field 'cvv' from Card is invalid, must have 3 digits")
	}
	c.DueDay = c.DueDayJson.Time
	c.ExpirationDate = c.ExpirationDateJson.Time
	//TODO: Check if expiration date is valid
	return nil
}
func (c *Card) validateUpdate() error {
	if c.AvailableCredit != 0 {
		if c.AvailableCredit > c.Limit {
			c.AvailableCredit = c.Limit
		} else if c.AvailableCredit < 0 {
			c.AvailableCredit = 0
		}
	}
	if c.DueDayJson != nil {
		c.DueDay = c.DueDayJson.Time
	}
	return nil
}

func (c *Card) fillNonDBFields() {
	c.DueDayJson = &dueDay{ Time: c.DueDay }
	c.ExpirationDateJson = &expirationDate{ Time: c.ExpirationDate }
}

//CardSorter type
type lessFunc func(p1, p2 *Card) bool

type cardSorter struct {
	cards []Card
	less  []lessFunc
}

func orderedBy(less ...lessFunc) *cardSorter {
	return &cardSorter{
		less: less,
	}
}

func (cs *cardSorter) Sort(cards []Card) {
	cs.cards = cards
	sort.Sort(cs)
}

func (cs *cardSorter) Len() int {
	return len(cs.cards)
}

func (cs *cardSorter) Swap(i, j int) {
	cs.cards[i], cs.cards[j] = cs.cards[j], cs.cards[i]
}

func (cs *cardSorter) Less(i, j int) bool {
	p, q := &cs.cards[i], &cs.cards[j]
	var k int
	for k = 0; k < len(cs.less)-1; k++ {
		less := cs.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return cs.less[k](p, q)
}