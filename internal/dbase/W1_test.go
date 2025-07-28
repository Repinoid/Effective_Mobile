package dbase

import (
	"emobile/internal/models"
	"time"
)

func (suite *TstHand) Test_01AddSubFunc() {

	// Gлобальный sub, template
	subG := models.Subscription{
		Service_name: "Yandex Plus",
		Price:        400,
		User_id:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		Start_date:   "01-02-2025",
		End_date:     "11-2025",
	}

	// трём таблицу передав пустую запись
	// впрочем, там и так нет ничего, на случай,
	// если в дальнейшем добавятся операции перед этим тестом
	_, err := suite.dataBase.DeleteSub(suite.ctx, models.Subscription{})
	suite.Require().NoError(err)

	// число подписок должно стать 0
	subs, err := suite.dataBase.ListSub(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(0, len(subs))

	cTag, err := suite.dataBase.AddSub(suite.ctx, subG)
	suite.Require().NoError(err)
	// EqualValues для унификации, т.к. RowsAffected() int64, а 1 - int
	suite.Require().EqualValues(1, cTag.RowsAffected())

	// структура "плоская", без указателей, поэтому копия независимая
	sub1 := subG
	sub1.Edt = time.Time{}
	sub1.Price = 0
	cTag, err = suite.dataBase.AddSub(suite.ctx, sub1)
	// ошибка т.к. такой же PRIMARY KEY (user_id, service_name)
	suite.Require().Error(err)
	// 0 - запись не вставилась Code = "23505"
	suite.Require().EqualValues(0, cTag.RowsAffected())
	suite.Require().Contains(err.Error(), "SQLSTATE 23505")

	// число подписок должно стать 2
	subs, err = suite.dataBase.ListSub(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().EqualValues(1, len(subs))

}
