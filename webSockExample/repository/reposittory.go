package repository

import (
	"context"
	"fmt"
	"github.com/Yideg/webSockExample/entity"
	"github.com/jackc/pgx/pgxpool"
)

type PgxRepo struct {
	pgx *pgxpool.Pool
	ctx context.Context
}

func NewPgx(p *pgxpool.Pool, c context.Context) *PgxRepo {
	return &PgxRepo{pgx: p, ctx: c}
}

//Add new Item the store
func (p PgxRepo) NewItem(pgx *pgxpool.Pool, ctx context.Context, new_item interface{}) {
	//aa,ok:=new_item.(Message)
	info := entity.Message{}
	mes := new_item.(map[string]interface{})
	for _, value := range mes {
		info.Text = value.(string)
	}
	fmt.Println("new message ", info)
	ProfileInsert := `INSERT INTO temp(message) VALUES($1);`
	_, err := pgx.Exec(ctx, ProfileInsert, info.Text)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("data inserted successfully")

}

//Update the existing  tem
func (p PgxRepo) UpdateItem(pgx *pgxpool.Pool, ctx context.Context, item interface{}) {
	ProfileUpdate := `UPDATE temp SET message =$1;`
	info := entity.Message{}
	mes := item.(map[string]interface{})
	for _, value := range mes {
		info.Text = value.(string)
	}
	_, err := pgx.Exec(ctx, ProfileUpdate, info.Text)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("data updated successfully")

}

//List all Items from the store
func (p PgxRepo) GetItems(pgx *pgxpool.Pool, ctx context.Context) []entity.Message {
	ProfileQuery := `SELECT * FROM temp;`
	rows, err := pgx.Query(ctx, ProfileQuery)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	var information []entity.Message
	for rows.Next() {
		info := entity.Message{}
		err = rows.Scan(&info.Text)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		information = append(information, info)
	}
	defer rows.Close()
	return information
}
