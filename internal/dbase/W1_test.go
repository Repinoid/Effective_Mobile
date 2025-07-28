package dbase

import (
	"emobile/internal/models"
	"time"

	"github.com/google/uuid"
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

	// генерируем user_id
	sub1.User_id = uuid.NewString()
	// повторяем добавление записи, с иным user_id
	cTag, err = suite.dataBase.AddSub(suite.ctx, sub1)
	suite.Require().NoError(err)
	// 1 - запись добавилась
	suite.Require().EqualValues(1, cTag.RowsAffected())

	// число подписок должно стать 2
	subs, err = suite.dataBase.ListSub(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().EqualValues(2, len(subs))

	sub1.Service_name = "Мурзилка"
	// повторяем добавление записи, с иным Service_name
	cTag, err = suite.dataBase.AddSub(suite.ctx, sub1)
	suite.Require().NoError(err)
	// 1 - запись добавилась
	suite.Require().EqualValues(1, cTag.RowsAffected())

	// число подписок должно стать 3
	subs, err = suite.dataBase.ListSub(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().EqualValues(3, len(subs))

	sub2 := sub1
	sub2.Service_name = "Völkischer Beobachter"
	// удаляем несуществующее ныне
	cTag, err = suite.dataBase.DeleteSub(suite.ctx, sub2)
	suite.Require().NoError(err)
	// 0 - нету такого, вот и не удалилось
	suite.Require().EqualValues(0, cTag.RowsAffected())

	// удаляем самую первую запись - subG
	cTag, err = suite.dataBase.DeleteSub(suite.ctx, subG)
	suite.Require().NoError(err)
	// 1 - норм
	suite.Require().EqualValues(1, cTag.RowsAffected())

	// 3-1 = 2
	subs, err = suite.dataBase.ListSub(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().EqualValues(2, len(subs))

	// подымем Мурзиле цену
	sub1.Price = 777
	cTag, err = suite.dataBase.UpdateSub(suite.ctx, sub1)
	suite.Require().NoError(err)
	// 1 - норм, апгрейд
	suite.Require().EqualValues(1, cTag.RowsAffected())

	// количество не изменилось
	subs, err = suite.dataBase.ListSub(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().EqualValues(2, len(subs))

	// так можно долго продолжать, надеюсь, достаточно тестов

}
