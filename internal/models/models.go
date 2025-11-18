package models

// TokenResponse representa a resposta da API de geração de token
type TokenResponse struct {
	Code        string `json:"code"`
	Message     string `json:"message,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
	DateTime    int64  `json:"datetime,omitempty"`
	Scope       string `json:"scope,omitempty"`
	Modulos     int    `json:"modulos,omitempty"`
}

// WSTestResponse representa a resposta da API de teste
type WSTestResponse struct {
	Code                  string `json:"code"`
	Message               string `json:"message,omitempty"`
	OrganizadorCodigo     int    `json:"organizadorCodigo,omitempty"`
	OrganizadorNome       string `json:"organizadorNome,omitempty"`
	OrganizadorCNPJ       string `json:"organizadorCnpj,omitempty"`
	OrganizadorLojaMatriz int    `json:"organizadorLojaMatriz,omitempty"`
	OrganizadorCodISGA    int    `json:"organizadorCodIsga,omitempty"`
	Versao                string `json:"versao,omitempty"`
	VersaoData            string `json:"versaoData,omitempty"`
	Erro                  string `json:"erro,omitempty"`
	Erro2                 string `json:"erro2,omitempty"`
}

// Application representa uma aplicação registrada no sistema
type Application struct {
	ClientID     string
	ClientSecret string
	JWTExpiracao int
	Scopo        int64
	Status       int
	Nome         string
}

// ZenviaWebhookRequest representa o payload recebido do webhook Zenvia para email
type ZenviaWebhookRequest struct {
	Type          string              `json:"type"`
	Message       ZenviaMessage       `json:"message"`
	MessageStatus ZenviaMessageStatus `json:"messageStatus"`
}

// ZenviaMessage representa a mensagem no webhook
type ZenviaMessage struct {
	To         string `json:"to"`
	ExternalID string `json:"externalId"`
	MessageId  string `json:"id"`
}

// ZenviaMessageStatus representa o status da mensagem
type ZenviaMessageStatus struct {
	Code        string              `json:"code"`
	Description string              `json:"description"`
	Causes      []ZenviaStatusCause `json:"causes"`
}

// ZenviaStatusCause representa a causa de um status
type ZenviaStatusCause struct {
	Reason  string `json:"reason"`
	Details string `json:"details"`
}

// ZenviaWebhookResponse representa a resposta do webhook
type ZenviaWebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
