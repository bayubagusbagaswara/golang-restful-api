package service

import (
	"context"
	"database/sql"
	"golangrestfulapi/helper"
	"golangrestfulapi/model/domain"
	"golangrestfulapi/model/web"
	"golangrestfulapi/repository"
)

type CategoryServiceImpl struct {

	// disini kita butuh repository, karena untuk manipulasi datanya melalui repository
	CategoryRepository repository.CategoryRepository
	// kita buth koneksi ke databasenya
	DB *sql.DB
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request web.CategoryCreateRequest) web.CategoryResponse {
	// karena kita menggunakan database transactional, maka untuk request kita butuh transactional
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)

	// pertama kita method diatas diwrap trancational, tapi juga saat terjadi error maka kita rollback
	defer helper.CommitOrRollback(tx)

	// create category
	category := domain.Category{
		Name: request.Name,
	}

	c := service.CategoryRepository.Save(ctx, tx, category)
	return helper.ToCategoryResponse(c)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request web.CategoryUpdateRequest) web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	// cek apakah ada category di database
	category, err := service.CategoryRepository.FindById(ctx, tx, request.Id)
	helper.PanicIfError(err)

	// jika ada data category, maka kita ubah data name category nya dengan data baru dari request
	category.Name = request.Name

	c := service.CategoryRepository.Update(ctx, tx, category)
	return helper.ToCategoryResponse(c)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)
	// cek apakah ada category di database
	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	helper.PanicIfError(err)

	// jika category nya ada di database, maka bisa kita hapus
	service.CategoryRepository.Delete(ctx, tx, category)
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int) web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	c, err2 := service.CategoryRepository.FindById(ctx, tx, categoryId)
	helper.PanicIfError(err2)

	return helper.ToCategoryResponse(c)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []web.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	categoryList := service.CategoryRepository.FindAll(ctx, tx)

	return helper.ToCategoryResponses(categoryList)
}
