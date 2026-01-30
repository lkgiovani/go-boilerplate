package delivery

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToResponseDTO(u *user.User) dto.UserResponseDTO {
	return dto.UserResponseDTO{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		ImgURL:     u.ImgURL,
		Admin:      u.Admin,
		Active:     u.Active,
		Source:     u.Source,
		Metadata:   m.ToMetadataDTO(u.Metadata),
		LastAccess: u.LastAccess,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (m *UserMapper) ToMetadataDTO(meta user.UserMetadata) dto.UserMetadataDTO {
	var proSource *string
	if meta.ProSource != nil {
		s := string(*meta.ProSource)
		proSource = &s
	}

	return dto.UserMetadataDTO{
		AccessMode:              string(meta.AccessMode),
		PlanType:                string(meta.PlanType),
		PlanExpirationDate:      meta.PlanExpirationDate,
		ProSource:               proSource,
		MaxResources:            meta.MaxResources,
		MaxRequestsPerMonth:     meta.MaxRequestsPerMonth,
		MaxAccounts:             meta.MaxAccounts,
		MaxCategoriesPerAccount: meta.MaxCategoriesPerAccount,
		MaxTransactionsPerMonth: meta.MaxTransactionsPerMonth,
		CanExportData:           meta.CanExportData,
		CanUseReports:           meta.CanUseReports,
		CanUseAdvancedFeatures:  meta.CanUseAdvancedFeatures,
		CanCreateBudgets:        meta.CanCreateBudgets,
		CanUseGoals:             meta.CanUseGoals,
		EmailVerified:           meta.EmailVerified,
		ReputationStatus:        string(meta.ReputationStatus),
		SuspiciousActivityCount: meta.SuspiciousActivityCount,
		LastSecurityCheck:       meta.LastSecurityCheck,
		LastPermissionCheck:     meta.LastPermissionCheck,
		Notes:                   meta.Notes,
		Locale:                  meta.Locale,
		Currency:                meta.Currency,
	}
}

func (m *UserMapper) ToPageResponse(users []user.User, total int64, page, size int) dto.PageResponse {
	content := make([]dto.UserResponseDTO, len(users))
	for i, u := range users {
		content[i] = m.ToResponseDTO(&u)
	}

	totalPages := int(total) / size
	if int(total)%size != 0 {
		totalPages++
	}

	return dto.PageResponse{
		Content:       content,
		TotalElements: total,
		TotalPages:    totalPages,
		Size:          size,
		Number:        page,
		First:         page == 1,
		Last:          page >= totalPages,
	}
}

func (h *Handler) SaveUser(c *fiber.Ctx) error {
	var req dto.UserPostRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	newUser := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Admin:    req.Admin != nil && *req.Admin,
		Active:   req.Active == nil || *req.Active,
		Source:   "LOCAL",
		Metadata: user.NewDefaultMetadata(),
	}
	newUser.Metadata.EmailVerified = true

	if req.Password != nil {
		newUser.Password = req.Password
	}

	if req.Source != nil {
		newUser.Source = *req.Source
	}

	if err := h.UserService.Repository.Create(c.Context(), newUser); err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusCreated).JSON(mapper.ToResponseDTO(newUser))
}

func (h *Handler) SignupUser(c *fiber.Ctx) error {
	var req dto.SignupUserRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	newUser := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Admin:    false,
		Active:   true,
		Source:   "LOCAL",
		Metadata: user.NewDefaultMetadata(),
	}

	newUser.Password = &req.Password

	if err := h.AuthService.Register(c.Context(), newUser); err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusCreated).JSON(mapper.ToResponseDTO(newUser))
}

func (h *Handler) ResendVerification(c *fiber.Ctx) error {
	var req dto.ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	if err := h.EmailVerificationService.ResendVerification(c.UserContext(), req.Email); err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Email de verificação reenviado com sucesso!",
	})
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	var req dto.UserPutRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	email := c.Locals("userEmail").(string)

	existingUser, err := h.UserService.Repository.GetByEmail(c.Context(), email)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	existingUser.Name = req.Name
	if existingUser.Email != req.Email {
		existingUser.Email = req.Email
		existingUser.Metadata.EmailVerified = false
	}
	if req.ImgURL != nil {
		existingUser.ImgURL = req.ImgURL
	}

	if err := h.UserService.Repository.Update(c.Context(), existingUser); err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(existingUser))
}

func (h *Handler) UpdatePassword(c *fiber.Ctx) error {
	var req dto.UserPutPasswordRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	email := c.Locals("userEmail").(string)

	existingUser, err := h.UserService.Repository.GetByEmail(c.Context(), email)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	currentPassword := ""
	if req.CurrentPassword != nil {
		currentPassword = *req.CurrentPassword
	}

	if err := h.UserService.Repository.ChangePassword(c.Context(), existingUser.ID, currentPassword, req.Password); err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) AddImage(c *fiber.Ctx) error {
	var req dto.UploadImageRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	result, err := h.StorageService.GetPresignedUploadUrl(c.Context(), req.FileName, req.ContentType, req.ContentLength)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.UploadResponseDTO{
		UploadSignedURL: result.SignedUrl,
		PublicURL:       result.FinalUrl,
	})
}

func (h *Handler) FindUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	foundUser, err := h.UserService.Repository.GetByID(c.Context(), id)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(foundUser))
}

func (h *Handler) FindUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	foundUser, err := h.UserService.Repository.GetByEmail(c.Context(), email)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(foundUser))
}

func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {

	email := c.Locals("userEmail").(string)

	foundUser, err := h.UserService.Repository.GetByEmail(c.Context(), email)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(foundUser))
}

func (h *Handler) DeleteUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	currentUserID := c.Locals("userID").(int64)
	if currentUserID == id {
		return errors.Errorf(errors.EINVALID, "You cannot delete your own account")
	}

	if err := h.UserService.Repository.Delete(c.Context(), id); err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) DeleteUsersByIDs(c *fiber.Ctx) error {
	idsParam := c.Query("ids")
	if idsParam == "" {
		return errors.Errorf(errors.EBADREQUEST, "IDs parameter is required")
	}

	// Parse IDs from format "1,2,3" or similar
	parts := strings.Split(idsParam, ",")
	var ids []int64
	for _, p := range parts {
		id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return errors.Errorf(errors.EBADREQUEST, "Valid IDs are required")
	}

	if err := h.UserService.Repository.DeleteByIDs(c.Context(), ids); err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) FindAllUsers(c *fiber.Ctx) error {
	keyword := c.Query("keyword", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	var users []user.User
	var total int64
	var err error

	if keyword != "" {
		users, total, err = h.UserService.Repository.FindAllWithFilter(c.Context(), keyword, page, size)
	} else {
		users, total, err = h.UserService.Repository.FindAll(c.Context(), page, size)
	}

	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	pageResponse := mapper.ToPageResponse(users, total, page, size)

	return c.Status(fiber.StatusOK).JSON(pageResponse)
}

func (h *Handler) ToggleUserStatus(c *fiber.Ctx) error {
	userIDParam := c.Params("userId")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	activeParam := c.Query("active")
	active := activeParam == "true"

	currentUserID := c.Locals("userID").(int64)
	if currentUserID == userID && !active {
		return errors.Errorf(errors.EINVALID, "You cannot deactivate your own account")
	}

	if err := h.UserService.Repository.ToggleStatus(c.Context(), userID, active); err != nil {
		return h.ErrorHandler(c, err)
	}

	updatedUser, err := h.UserService.Repository.GetByID(c.Context(), userID)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}

func (h *Handler) UpdateUserAdmin(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	var req dto.UserPutRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	existingUser, err := h.UserService.Repository.GetByID(c.Context(), id)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	existingUser.Name = req.Name
	if existingUser.Email != req.Email {
		existingUser.Email = req.Email
		existingUser.Metadata.EmailVerified = false
	}
	if req.ImgURL != nil {
		existingUser.ImgURL = req.ImgURL
	}

	if err := h.UserService.Repository.Update(c.Context(), existingUser); err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(existingUser))
}

func (h *Handler) UpdatePasswordAdmin(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	var req dto.UserPutPasswordRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	if err := h.UserService.Repository.ResetUserPassword(c.Context(), id, req.Password); err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) UpdateAccessMode(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	var req dto.UserUpdateAccessModeDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	updatedUser, err := h.UserService.Repository.UpdateAccessMode(c.Context(), id, req.AccessMode)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}

func (h *Handler) UpdateFeatures(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	var req dto.UserUpdateFeaturesDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	updatedUser, err := h.UserService.Repository.UpdateFeatures(c.Context(), id, req.CanCreateBudgets, req.CanExportData, req.CanUseReports, req.CanUseGoals)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}

func (h *Handler) UpdateLimits(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	var req dto.UserUpdateLimitsDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	updatedUser, err := h.UserService.Repository.UpdateLimits(c.Context(), id, req.MaxAccounts, req.MaxTransactionsPerMonth, req.MaxCategoriesPerAccount)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}

func (h *Handler) GrantLifetimePro(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	var req dto.UserGrantLifetimeProDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	updatedUser, err := h.UserService.Repository.GrantLifetimePro(c.Context(), id, req.Reason)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}

func (h *Handler) EnsureMetadata(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	updatedUser, err := h.UserService.Repository.EnsureMetadata(c.Context(), id)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}

func (h *Handler) RevokeLifetimePro(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid user ID")
	}

	currentUserID := c.Locals("userID").(int64)
	if currentUserID == id {
		return errors.Errorf(errors.EINVALID, "You cannot revoke your own Lifetime Pro status")
	}

	updatedUser, err := h.UserService.Repository.RevokeLifetimePro(c.Context(), id)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusOK).JSON(mapper.ToResponseDTO(updatedUser))
}
