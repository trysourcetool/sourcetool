package adapters

import (
	"github.com/trysourcetool/sourcetool/backend/dto"
	"github.com/trysourcetool/sourcetool/backend/server/http/responses"
)

func DeleteUserOutputToResponse(out *dto.DeleteUserOutput) *responses.DeleteUserResponse {
	if out == nil {
		return nil
	}

	return &responses.DeleteUserResponse{
		Success: out.Success,
	}
}
