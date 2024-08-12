package server

import (
	"Marketplace/internal/entities"
	"errors"
	"fmt"
	"net/http"
)

func (ad *Server) checkAdOwner(r *http.Request, adID int) (error, *entities.Advert) {
	//curUserId := r.Context().Value("user_id")
	//if curUserId == nil {
	//	return errors.New("User not logged in"), nil
	//}
	curUserId := 4
	curAdv, err := ad.db.GetAdvert(r.Context(), adID)
	if err != nil {
		return fmt.Errorf("unable too get advert: %w", err), nil
	}
	if curUserId != curAdv.UserId {
		return errors.New("the ad doesn't belong to this user"), nil
	}
	return nil, curAdv
}
