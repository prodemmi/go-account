package handlers

import (
	"net/http"

	"aman/internal/ports/inbound"

	"github.com/gin-gonic/gin"
)

// SessionHandler exposes the current user's sessions. CAS does not manage
// users: a session belongs to the app + identifier pair authenticated by the
// request middleware.
type SessionHandler struct {
	sessions inbound.SessionUseCase
}

func NewSessionHandler(sessions inbound.SessionUseCase) *SessionHandler {
	return &SessionHandler{sessions: sessions}
}

// Index lists the caller's active sessions.
// @Summary List my sessions (CAS)
// @Tags Sessions
// @Produce json
// @Param APP-ID header string true "app id"
// @Success 200 {object} map[string]interface{}
// @Router /v1/user/sessions [get]
func (h SessionHandler) Index(c *gin.Context) {
	sess, ok := currentSession(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	list, err := h.sessions.List(c.Request.Context(), sess.AppID.String(), sess.Identifier)
	if err != nil {
		redirectNotFound(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   list,
	})
}

// Delete terminates one session by id.
// @Summary Terminate a session (CAS)
// @Tags Sessions
// @Param id path string true "session id"
// @Success 200 {object} map[string]interface{}
// @Router /v1/user/sessions/{id} [delete]
func (h SessionHandler) Delete(c *gin.Context) {
	sess, ok := currentSession(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	if err := h.sessions.Terminate(c.Request.Context(), sess.AppID.String(), sess.Identifier, c.Param("id")); err != nil {
		redirectNotFound(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": true})
}

// DeleteBatch terminates several sessions at once, or all of them when the
// request has no body / no ids.
// @Summary Terminate sessions in batch (CAS)
// @Tags Sessions
// @Param payload body terminateSessionsBody false "session ids"
// @Success 200 {object} map[string]interface{}
// @Router /v1/user/sessions [delete]
func (h SessionHandler) DeleteBatch(c *gin.Context) {
	sess, ok := currentSession(c)
	if !ok {
		redirectNotFound(c)
		return
	}

	var payload terminateSessionsBody
	_ = c.ShouldBindJSON(&payload)

	ctx := c.Request.Context()
	appID := sess.AppID.String()

	if len(payload.IDs) > 0 {
		if err := h.sessions.TerminateMany(ctx, appID, sess.Identifier, payload.IDs); err != nil {
			redirectNotFound(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"terminated": len(payload.IDs)}})
		return
	}

	if err := h.sessions.TerminateAll(ctx, appID, sess.Identifier); err != nil {
		redirectNotFound(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"terminated": "all"}})
}
