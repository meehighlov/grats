package errs

import "errors"

var ErrSaveWishValidation = errors.New("add save wish validation error")
var ErrEditWishNameValidation = errors.New("edit wish name validation error")
var ErrSaveEditLinkValidation = errors.New("edit link validation error")
var ErrEditPriceValidation = errors.New("edit price validation error")
