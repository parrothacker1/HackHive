package middlewares

import (
	"context"
	"net/http"

	"github.com/parrothacker1/Solvelt/users/models"
	"github.com/parrothacker1/Solvelt/users/utils/database"
	"gorm.io/gorm"
)

var TeamLeaderMiddleware http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
  user_id := r.Context().Value("user_id")
  team_id := r.Context().Value("team_id")
  var team models.Team
  var ctx context.Context
  if err := database.DB.Table("teams").Where("team_id = ? and team_leader = ?",team_id,user_id).First(&team).Error;err != nil {
    if err == gorm.ErrRecordNotFound {
      ctx = context.WithValue(r.Context(),"stop",true)
      r = r.WithContext(ctx)
      http.Error(w,`{"status":"fail","message":"Team with this leader does not exists"}`,http.StatusNotFound)
      return
    } else {
      ctx = context.WithValue(r.Context(),"stop",true)
      r = r.WithContext(ctx)
      http.Error(w,`{"status":"error","message":"Error in fetching data"}`,http.StatusInternalServerError)
      return
    }
  }
  ctx = context.WithValue(r.Context(),"leader",true)
  r = r.WithContext(ctx)
}
