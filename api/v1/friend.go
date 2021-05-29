package apiv1

import (
	"container/list"
	"fmt"
	. "litrocket/common"
	"litrocket/model"
	"litrocket/utils/dataencry"
)

type Friend struct {
	Url    string
	SrcID  UserID
	DestID UserID
	Notes  string
}

type FriResult struct {
	Url      string
	Code     int
	FriendID UserID
}

// 邻接表
type Graph struct {
	Vertexs int          // 顶点数
	Ids     map[int]int  // 保存每个顶点代表的用户ID
	List    []*list.List // 链表数组,每个顶点一个链表,保存与其有关系的顶点
}

// 初始化图
func (G *Graph) InitGraph(vertexs int) {
	G.Vertexs = vertexs
	G.Ids = make(map[int]int, vertexs)
	G.List = make([]*list.List, vertexs)
	for i := 0; i < vertexs; i++ {
		G.List[i] = list.New()
	}
}

// 添加边
func (G *Graph) Add(start, end int) {
	G.List[start].PushBack(end)
	G.List[end].PushBack(start)
}

// 广度优先遍历
// 类似层序遍历
// 除自身和一度好友外,找到所有二度好友,将两次或多次出现的二度好友推荐给用户.
func (G *Graph) bfs(start int) map[int]int {
	var (
		queue    Queue
		visisted = make([]bool, G.Vertexs)
		more     = make(map[int]int, G.Vertexs)
		onelevel = list.New() //一度好友序列
	)

	for i := 0; i < G.Vertexs; i++ {
		visisted[i] = false
	}

	visisted[start] = true
	queue.InitQueue()
	queue.Enqueue(start)

	for queue.Size() != 0 {
		w := queue.Head()
		queue.Dequeue()
		// fmt.Printf("%d ", w)  这里输出即为广度遍历的序列
		for e := G.List[w].Front(); e != nil; e = e.Next() {
			q := e.Value.(int)

			// 记录一度好友
			if w == start {
				onelevel.PushBack(q)
			}

			// 把q放入more这个哈希表中,并记录出现的次数
			if v, ok := more[q]; ok {
				v++
				more[q] = v
			} else {
				more[q] = 1 // 第一次出现
			}

			if !visisted[q] {
				visisted[q] = true
				queue.Enqueue(q)
			}
		}
	}

	// more 中去掉顶点和一度好友,得到的就是二度好友,以及每个二度好友出现的次数
	fmt.Println(more)

	delete(more, start)
	for e := onelevel.Front(); e != nil; e = e.Next() {
		delete(more, e.Value.(int))
	}

	fmt.Println(more)

	return more
}

// Add a friend to table "friend", friend's state is "waiting friend agree".
func AddFriend(json []byte) {
	var (
		err    error
		friend Friend

		result struct {
			Url    string
			Code   int
			Friend model.User
		}

		mess struct {
			Url    string
			Friend model.User
		}
	)

	if err = dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	// 查询用户ID
	user, exist := model.SearchUser(friend.Notes)
	if !exist {
		return
	}

	result.Url = "add/friresult"
	result.Friend, result.Code = model.AddFriend(friend.SrcID, UserID(user.ID))

	// 被加的用户在线,发送消息
	if val, ok := AllUsers.Load(UserID(user.ID)); ok { //! 大坑, 如果不把ID转为UserID类型,就会检测出不在线,查询时不仅KEY的值要一样,而且KEY的类型也要和存这个键值对的时候一样
		mess.Url = friend.Url
		mess.Friend, _ = model.SearchByID(friend.SrcID)
		conns := val.(Conns)
		r, _ := dataencry.Marshal(mess)
		b := append(r, []byte("\r\n--\r\n")...)
		conns.ResponseConn.Write(b)
	}

	// 返回加好友的结果
	if val, ok := AllUsers.Load(friend.SrcID); ok {
		conns := val.(Conns)
		b, _ := dataencry.Marshal(result)
		r := append(b, []byte("\r\n--\r\n")...)
		conns.ResponseConn.Write(r)
	}
}

func GetAllFriend(json []byte) {
	var (
		friend  Friend
		friends []model.Friend
		result  struct {
			Url    string
			Friend []model.Friend
		}
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	friends = model.GetAllFriend(friend.SrcID)

	result.Friend = friends
	result.Url = friend.Url
	buf, _ := dataencry.Marshal(result)
	buf = append(buf, []byte("\r\n--\r\n")...)

	if conns, ok := AllUsers.Load(friend.SrcID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(buf)
	}
}

func GetFriendInfo(json []byte) {
	var (
		exist  bool
		friend struct {
			Url      string
			SrcID    UserID
			FriendID UserID
		}
	)

	var (
		result struct {
			Url  string
			Info model.User
		}
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	result.Url = friend.Url
	result.Info, exist = model.SearchByID(friend.FriendID)
	if exist {
		if conns, ok := AllUsers.Load(friend.SrcID); ok {
			conn := conns.(Conns)
			buf, _ := dataencry.Marshal(result)
			buf = append(buf, []byte("\r\n--\r\n")...)
			conn.ResponseConn.Write(buf)
		}
	}

}

func DelFriend(json []byte) {
	var (
		friend Friend
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	model.DelFriend(friend.SrcID, friend.DestID)

	// 目标在线则通知,不在线不通知
	if conns, ok := AllUsers.Load(friend.DestID); ok {
		conn := conns.(Conns)
		b := append(json, []byte("\r\n--\r\n")...)
		conn.ResponseConn.Write(b)
	}
}

func EditFriendNotes(json []byte) {
	var (
		friend Friend
		result FriResult
	)

	if err := dataencry.Unmarshal(json, &friend); err != nil {
		return
	}

	i := model.EditFriendNotes(friend.SrcID, friend.DestID, friend.Notes)
	result.Code = i
	result.FriendID = friend.DestID
	result.Url = friend.Url
	r, _ := dataencry.Marshal(result)

	if conns, ok := AllUsers.Load(friend.SrcID); ok {
		conn := conns.(Conns)
		conn.ResponseConn.Write(r)
	}
}

// Friend Recommand
// THink: A有B,C两个好友,B和C有共同好友D,那么将D推荐给A, 即如果有两个或多个一度好友有同样的二度好友(除自己),即推荐这个同样的二度好友.  也就是共同好友(多数都是这样)
// How To Do :
// 	(1) 将用户根据其好友关系组成一个图
// 	(2) 以此用户为顶点,寻找一度好友(我的好友)和每一个一度好友的二度好友(好友的好友)
//	(3) 查询数据库一度好友,添加至邻接表
func FriendRecommand(json []byte) {
	var (
		ReComd struct {
			Url string
			Id  UserID
		}

		result struct {
			Url  string
			Name []string
		}

		G Graph
	)

	if err := dataencry.Unmarshal(json, &ReComd); err != nil {
		return
	}

	// 根据好友关系创建图
	G.InitGraph(50)
	CreateFriendGraph(&G, ReComd.Id)

	// DEBUG.
	// for i := 0; i < 3; i++ {
	// 	fmt.Printf("[%d]: ", i)
	// 	for e := G.List[i].Front(); e != nil; e = e.Next() {
	// 		fmt.Printf("%d->", e.Value.(int))
	// 	}
	// 	fmt.Println()
	// }

	if v, ok := G.Ids[(int(ReComd.Id))]; ok {
		mayfriends := G.bfs(v)

		// 遍历,将好友信息取出返回给客户端
		// todo if len(mayfriends) > 10  // 排一下序,将最适合的返回
		i := 0
		for k, _ := range mayfriends {
			result.Url = ReComd.Url
			result.Name = make([]string, len(mayfriends))
			for id, index := range G.Ids {
				if k != index {
					continue
				}

				user, _ := model.SearchByID(UserID(id))
				result.Name[i] = user.UserName
				i++
			}
		}

		b, _ := dataencry.Marshal(result)
		r := append(b, []byte("\r\n--\r\n")...)

		if conns, ok := AllUsers.Load(ReComd.Id); ok {
			conn := conns.(Conns)
			conn.ResponseConn.Write(r)
		}
	}
}

// 建图只需要查到好友的好友即可,不需继续往下查找建图
func CreateFriendGraph(G *Graph, Id UserID) {
	index := 0
	friends := model.GetAllFriend(Id)
	G.Ids[int(Id)] = index

	// 先把一度好友存入,顶点与一度好友建立关系
	for i := 0; i < len(friends); i++ {
		index++
		G.Ids[int(friends[i].FriendID)] = index
		G.Add(0, index)
	}

	// 再把二度好友存入,让二度好友与一度好友建立关系,二度好友不会与顶点建立关系,因为与顶点有关系的都是一度好友,若某组二度好友中有顶点,则不建立关系
	for i := 0; i < len(friends); i++ {
		friend := model.GetAllFriend(friends[i].FriendID)
		v, _ := G.Ids[int(friends[i].FriendID)]

		for j := 0; j < len(friend); j++ {
			// 二度好友不与顶点建立关系,只与一度好友建立关系
			if friend[j].FriendID == Id {
				continue
			}

			// 如果此二度好友已经存在于Ids数组中,说明已经与某一度用户建立了关系
			if v2, ok := G.Ids[int(friend[j].FriendID)]; ok {
				// 如果此二度好友已经与这个一度好友建立了关系,那么不再次建立关系
				if IsExist(G.List[v], v2) {
					continue
				}

				// 未与这个一度好友建立关系,那么建立关系
				G.Add(v, v2)
				continue
			}

			// 如果不存在Ids中,说明未与任何一度用户建立关系
			index++
			G.Ids[int(friend[j].FriendID)] = index
			G.Add(v, index)
		}
	}
}

func IsExist(list *list.List, val int) bool {
	for e := list.Front(); e != nil; e = e.Next() {
		if val == e.Value.(int) {
			return true
		}
	}

	return false
}
