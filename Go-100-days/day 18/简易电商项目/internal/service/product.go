// internal/service/product.go
// 商品业务逻辑层
// 学习要点：商品 CRUD 业务逻辑，调用 DAO 层
//
// 业务逻辑相对简单，主要是参数转换和调用 DAO
// 复杂业务（如促销、折扣）可在此层扩展

package service

import (
	"errors"

	"gorm.io/gorm"

	"simple-ecommerce/internal/dao"
	"simple-ecommerce/internal/model"
)

// ProductService 商品业务服务
type ProductService struct {
	productDAO *dao.ProductDAO
}

// NewProductService 创建 ProductService 实例
func NewProductService() *ProductService {
	return &ProductService{
		productDAO: dao.NewProductDAO(),
	}
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(product *model.Product) error {
	if product.Price < 0 {
		return errors.New("价格不能为负数")
	}
	if product.Stock < 0 {
		return errors.New("库存不能为负数")
	}
	return s.productDAO.Create(product)
}

// GetProduct 查询商品详情
func (s *ProductService) GetProduct(id uint) (*model.Product, error) {
	product, err := s.productDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商品不存在")
		}
		return nil, err
	}
	return product, nil
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(product *model.Product) error {
	// 先检查商品是否存在
	_, err := s.productDAO.GetByID(product.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在")
		}
		return err
	}
	return s.productDAO.Update(product)
}

// DeleteProduct 删除商品
func (s *ProductService) DeleteProduct(id uint) error {
	return s.productDAO.Delete(id)
}

// ProductListResult 商品列表查询结果
type ProductListResult struct {
	List  []model.Product `json:"list"`  // 商品列表
	Total int64           `json:"total"` // 总数
	Page  int             `json:"page"`  // 当前页码
	Size  int             `json:"size"`  // 每页条数
}

// GetProductList 分页查询商品列表
// 学习要点：对分页参数做默认值处理，防止传入非法值
func (s *ProductService) GetProductList(page, size int) (*ProductListResult, error) {
	// 参数校验与默认值
	if page < 1 {
		page = 1 // 默认第1页
	}
	if size < 1 || size > 100 {
		size = 10 // 默认每页10条，最大100条
	}

	products, total, err := s.productDAO.GetList(page, size)
	if err != nil {
		return nil, err
	}

	return &ProductListResult{
		List:  products,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}
