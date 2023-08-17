package dao

import (
	"database/sql"
	"fmt"
	"spiderGo/db"
	"sync"
)

var (
	f               *fuDaNewsDao
	fuDaNewsDaoOnce sync.Once
)

type fuDaNewsDao struct {
	db *sql.DB
}

func (f *fuDaNewsDao) Insert(fuDaNewsList [][]string) (total, success int, err error) {
	//projects := []struct{ mascot string     release int }{{"tux", 1991}, {"duke", 1996}, {"gopher", 2009}, {"moby dock", 2013}}
	//data = [[发布日期:  2023-08-16 作者： 发展规划与学科建设处 福州大学世界排名再攀升 2023软科世界大学学术排名首次进入全球前300名] [发布日期:  2023-08-15 作者： 党委宣传部 【主题教育】福州大学召开主题教育测评会与座谈会]
	total = len(fuDaNewsList)
	stmt, err := f.db.Prepare("INSERT INTO fuDaNews(pubtime, author, title) VALUES( ?, ?, ?)")
	if err != nil {
		err = fmt.Errorf("预插入数据库失败:%w",err)
		return
	}
	defer stmt.Close()
	for _, fuDaNew := range fuDaNewsList {
		rs, err := stmt.Exec(fuDaNew[0], fuDaNew[1], fuDaNew[2])
		if err != nil {
			err = fmt.Errorf("插入数据库失败:%w",err)
			return total, success, err
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			err = fmt.Errorf("插入数据库后获取总数失败:%w",err)
			continue
		}
		success += int(affected)
	}
	return
}

func newFuDaNewsDao() *fuDaNewsDao {
	return &fuDaNewsDao{
		db:db.GetDbSingleInstance(),
	}
}

type FuDaNewsDaoInstance interface {
	Insert([][]string) (int, int, error)
}

func GetFuDaNewsDao() FuDaNewsDaoInstance {
	fuDaNewsDaoOnce.Do(func() {
		//logger.Logger.Println("start newFuDaNewsDao")
		f = newFuDaNewsDao()
		//logger.Logger.Printf("newFuDaNewsDao:%+v", f)
	})
	return f
}
