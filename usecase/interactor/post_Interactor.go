package interactor

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"

	"github.com/yzmw1213/PostService/db"
	"github.com/yzmw1213/PostService/domain/model"
	"github.com/yzmw1213/PostService/grpc/userservice"
	"github.com/yzmw1213/PostService/usecase/repository"
)

var (
	err       error
	post      model.Post
	posts     []model.Post
	postTag   model.PostTag
	postTags  []model.PostTag
	likeUser  model.PostLikeUser
	likeUsers []model.PostLikeUser
	rows      *sql.Rows
	validate  *validator.Validate
)

// PostInteractor 投稿サービスを提供するメソッド群
type PostInteractor struct{}

var _ repository.PostRepository = (*PostInteractor)(nil)

// Create 投稿1件を作成
func (p *PostInteractor) Create(postData *model.JoinPost) (*model.JoinPost, error) {
	validate = validator.New()
	post := postData.Post
	tags := postData.PostTags

	// Post構造体のバリデーション
	if err := validate.Struct(post); err != nil {
		return postData, err
	}

	// トランザクション開始
	tx := db.StartBegin()

	// 投稿登録
	if err := tx.Create(post).Error; err != nil {
		db.EndRollback()
		return postData, err
	}

	postID := post.ID
	// 投稿とタグ紐付け情報登録
	for _, tag := range tags {
		tag.PostID = postID
		if err := tx.Create(tag).Error; err != nil {
			db.EndRollback()
			return postData, err
		}
	}
	// トランザクションを終了しコミット
	db.EndCommit()
	return postData, err
}

// DeleteByID 指定されたIDに対する投稿1件を削除
func (p *PostInteractor) DeleteByID(id uint32) error {
	var post model.Post
	var postTag model.PostTag

	// トランザクション開始
	tx := db.StartBegin()
	// 指定されたPostIDのPostを削除
	if err := tx.Where("id = ?", id).Delete(&post).Error; err != nil {
		db.EndRollback()
		return err
	}
	// 指定されたPostIDのPostTagを削除
	if err := tx.Where("post_id = ?", id).Delete(&postTag).Error; err != nil {
		db.EndRollback()
		return err
	}
	// トランザクションを終了しコミット
	db.EndCommit()
	return nil
}

// List 投稿を全件取得
func (p *PostInteractor) List() ([]model.JoinPost, error) {
	rows, err := listAll(context.Background())
	if err != nil {
		fmt.Println("Error happened")
		return []model.JoinPost{}, err
	}
	// 取得したpostsに紐付け情報を付与して返す
	return createJoinPosts(rows)
}

// listAll 全件取得
func listAll(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post
	DB := db.GetDB()

	_, err := DB.Find(&posts).Rows()
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}
	return posts, nil
}

// Update 投稿を更新する
func (p *PostInteractor) Update(postData *model.JoinPost) (*model.JoinPost, error) {
	validate = validator.New()
	post := postData.Post
	tags := postData.PostTags

	// Post構造体のバリデーション
	if err := validate.Struct(post); err != nil {
		return postData, err
	}

	// トランザクション開始
	tx := db.StartBegin()

	if err := tx.Model(&post).Update(&postData.Post).Error; err != nil {
		db.EndRollback()
		return postData, err
	}

	// 投稿とタグ紐付け情報を全て削除
	deletePostTagByPostID(post.ID)

	postID := post.ID
	// 投稿とタグ紐付け情報登録
	for _, tag := range tags {
		tag.PostID = postID
		if err := tx.Create(tag).Error; err != nil {
			db.EndRollback()
			return postData, err
		}
	}
	// トランザクションを終了しコミット
	db.EndCommit()
	return postData, nil
}

// GetByID IDを元に投稿を1件取得する
func (p *PostInteractor) GetByID(ID uint32) (model.Post, error) {
	DB := db.GetDB()
	row := DB.First(&post, ID)
	if err := row.Error; err != nil {
		log.Printf("Error happend while Read for ID: %v\n", ID)
		return model.Post{}, err
	}
	DB.Table(db.PostTableName).Scan(row)
	return post, nil
}

// GetJoinPostByID IDを元に投稿、紐付け情報を1件取得する
func (p *PostInteractor) GetJoinPostByID(ID uint32) (model.JoinPost, error) {
	post, err := p.GetByID(ID)

	if err != nil {
		log.Printf("Error happend while Read for ID: %v\n", ID)
		return model.JoinPost{}, err
	}

	joinPost, err := createJoinPostSingle(post)
	if err != nil {
		log.Printf("Error happend while Read for ID: %v\n", ID)
		return model.JoinPost{}, err
	}

	return joinPost, nil
}

// Like 投稿のお気に入り
func (p *PostInteractor) Like(postData *model.PostLikeUser) (*model.PostLikeUser, error) {
	DB := db.GetDB()
	if err := DB.Create(postData).Error; err != nil {
		return postData, err
	}
	return postData, nil
}

// NotLike 投稿のお気に入りの取り消し
func (p *PostInteractor) NotLike(postData *model.PostLikeUser) (*model.PostLikeUser, error) {
	DB := db.GetDB()
	if err := DB.Where("post_id = ? ", postData.PostID).Where("user_id = ? ", postData.UserID).Delete(postData).Error; err != nil {
		return postData, err
	}
	return postData, nil
}

// listPostTagsByID PostIDを元にpostTagを検索し返す
func listPostTagsByID(ID uint32) ([]model.PostTag, error) {
	var postTagList []model.PostTag
	DB := db.GetDB()
	rows, err := DB.Where("post_id = ?", ID).Find(&postTags).Rows()
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}
	for rows.Next() {
		DB.ScanRows(rows, &postTag)
		postTagList = append(postTagList, postTag)
	}

	return postTagList, nil
}

// listPostLikeUsersByID PostIDを元にお気に入りしているユーザーを検索し返す
func listPostLikeUsersByID(ID uint32) ([]model.PostLikeUser, error) {
	var postLikeUserList []model.PostLikeUser
	DB := db.GetDB()
	rows, err := DB.Where("post_id = ?", ID).Find(&likeUsers).Rows()
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}
	for rows.Next() {
		DB.ScanRows(rows, &likeUser)
		postLikeUserList = append(postLikeUserList, likeUser)
	}

	return postLikeUserList, nil
}

func createJoinPosts(posts []model.Post) ([]model.JoinPost, error) {
	var joinPosts []model.JoinPost
	// 全ユーザーをUserServiceから取得
	users := getUserData()
	if err != nil {
		return []model.JoinPost{}, err
	}
	for _, post := range posts {
		var likeUsers []model.User
		// 投稿者のユーザー情報
		createUser := users[post.CreateUserID]

		// 紐付けられているタグ情報
		// タグ名はフロント側で保持しているタグストアから取得する
		postTags, err := listPostTagsByID(post.ID)
		if err != nil {
			return []model.JoinPost{}, err
		}
		// Likeしているユーザー
		// post.IDより取得
		postLikeUsers, err := listPostLikeUsersByID(post.ID)

		for _, user := range postLikeUsers {
			likeUsers = append(likeUsers, users[user.UserID])
		}
		// post.IDより取得
		joinPost := makeJoinPost(post, createUser, postTags, likeUsers)
		joinPosts = append(joinPosts, joinPost)
	}

	return joinPosts, nil
}

// 単一投稿のJoinPostを返す
func createJoinPostSingle(post model.Post) (model.JoinPost, error) {
	posts := []model.Post{post}
	joinPost, err := createJoinPosts(posts)

	return joinPost[0], err
}

func makeJoinPost(post model.Post, createUser model.User, postTags []model.PostTag, likeUsers []model.User) model.JoinPost {
	return model.JoinPost{
		Post:          &post,
		User:          &createUser,
		PostTags:      postTags,
		PostLikeUsers: likeUsers,
	}
}

func countPostTag() int {
	var count int
	DB := db.GetDB()
	DB.Model(&postTag).Count(&count)
	return count
}

// countPostTagByPostID IDを元に投稿に付けられているタグの件数を取得する
func countPostTagByPostID(ID uint32) int {
	var count int
	DB := db.GetDB()
	DB.Where("post_id = ?", ID).Model(&postTag).Count(&count)

	return count
}

// countPostLikeUserByPostID IDを元に投稿にお気に入りしているユーザー数を取得する
func countPostLikeUserByPostID(ID uint32) int {
	var count int
	DB := db.GetDB()
	DB.Where("post_id = ?", ID).Model(&likeUser).Count(&count)

	return count
}

func deletePostTagByPostID(ID uint32) {
	DB := db.GetDB()
	DB.Where("post_id = ?", ID).Delete(&postTags)
}

// ユーザーサービスからユーザー情報取得
func getUserData() map[uint32]model.User {
	proxyServerURL := os.Getenv("PROXY_SERVER_URL")

	cc, err := grpc.Dial(proxyServerURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()
	userClient := userservice.NewUserServiceClient(cc)

	request := &userservice.ListUserRequest{}
	res, err := userClient.ListUser(context.Background(), request)
	if err != nil {
		panic(err)
	}

	resUsers := res.GetProfile()
	var users map[uint32]model.User
	users = map[uint32]model.User{}
	for _, user := range resUsers {
		id := user.UserId

		users[id] = model.User{
			ID:       user.UserId,
			UserName: user.UserName,
		}
	}
	return users

}
