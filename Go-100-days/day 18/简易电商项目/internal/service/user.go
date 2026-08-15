// internal/service/user.go
// 用户业务逻辑层
// 学习要点：Service 层处理业务规则，调用 DAO 层操作数据库
//
// 为什么要有 Service 层（而不直接在 controller 操作 DAO）：
// 1. 业务逻辑复用：同一个业务可能被多个 controller 调用
// 2. 职责分离：controller 处理 HTTP，service 处理业务
// 3. 易测试：service 不依赖 HTTP 上下文，好写单元测试
//
// 注册流程：校验用户名是否已存在 -> 密码加密 -> 保存数据库
// 登录流程：查用户 -> 校验密码 -> 生成 JWT token

package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"simple-ecommerce/internal/dao"
	"simple-ecommerce/internal/model"
	"simple-ecommerce/pkg/jwt"
)

// UserService 用户业务服务
type UserService struct {
	userDAO *dao.UserDAO
}

// NewUserService 创建 UserService 实例
// 学习要点：依赖注入，service 持有 dao 的引用
func NewUserService() *UserService {
	return &UserService{
		userDAO: dao.NewUserDAO(),
	}
}

// Register 用户注册
// 学习要点：密码加密用 bcrypt，而非 md5
//
// 为什么用 bcrypt 不用 md5：
// 1. md5 已被破解，彩虹表可反查
// 2. bcrypt 自带盐值（salt），相同密码加密结果不同
// 3. bcrypt 可调节计算成本（cost），抵抗暴力破解
// 4. Go 标准库提供，无需第三方依赖
func (s *UserService) Register(username, password, phone string) (*model.User, error) {
	// 1. 检查用户名是否已存在
	_, err := s.userDAO.GetUserByUsername(username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}
	// 学习要点：gorm.ErrRecordNotFound 表示没找到，是预期的（说明用户名可用）
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("查询用户失败: " + err.Error())
	}

	// 2. 密码加密
	// 学习要点：bcrypt.GenerateFromPassword 生成加密后的密码
	// bcrypt.DefaultCost 是计算成本（10），越大越安全但越慢
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败: " + err.Error())
	}

	// 3. 创建用户
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Phone:    phone,
	}
	if err := s.userDAO.CreateUser(user); err != nil {
		return nil, errors.New("创建用户失败: " + err.Error())
	}

	return user, nil
}

// Login 用户登录
// 学习要点：登录返回 JWT token，前端保存后每次请求带上
func (s *UserService) Login(username, password string) (string, error) {
	// 1. 查询用户
	user, err := s.userDAO.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("用户不存在")
		}
		return "", errors.New("查询用户失败: " + err.Error())
	}

	// 2. 校验密码
	// 学习要点：bcrypt.CompareHashAndPassword 比较明文和密文
	// 返回 nil 表示匹配，返回 error 表示不匹配
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("密码错误")
	}

	// 3. 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", errors.New("生成 token 失败: " + err.Error())
	}

	return token, nil
}
