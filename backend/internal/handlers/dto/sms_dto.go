package dto

type SendSMSRequest struct {
	Phone    string  `json:"phone" validate:"required"`
	Content  string  `json:"content" validate:"required,max=1600"`
	DeviceID *string `json:"device_id"`
	SimSlot  int     `json:"sim_slot"`
}

type SendSMSResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	DeviceID  string `json:"device_id"`
}

type DeviceStatusDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Status     string  `json:"status"`
	Battery    int     `json:"battery"`
	LastSeenAt *string `json:"last_seen_at,omitempty"`
}

type DevicesStatusResponse struct {
	Devices []DeviceStatusDTO `json:"devices"`
}
