package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const (
	ViaCEP_URL = "https://viacep.com.br/ws/%s/json/"
)

type CEPService interface {
	GetAddressByCEP(cep string) (*ViaCEPResponse, error)
}

// ViaCEPService is a service to interact with the ViaCEP API
type ViaCEPService struct {
	BaseHttpService
}

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        bool   `json:"erro,omitempty"`
}

// NewViaCEPService creates a new ViaCEPService
func NewViaCEPService() CEPService {
	return &ViaCEPService{BaseHttpService{Client: &http.Client{}}}
}

// GetAddressByCEP returns the address for a given CEP
func (v *ViaCEPService) GetAddressByCEP(cep string) (*ViaCEPResponse, error) {

	resp, err := v.Client.Get(fmt.Sprintf(ViaCEP_URL, cep))
	if err != nil {
		slog.Error("Error getting address by CEP: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body: ", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, ErrInvalidCEP
	}

	var viaCepResponse ViaCEPResponse
	err = json.Unmarshal(body, &viaCepResponse)
	if err != nil {
		slog.Error("Error unmarshalling response body: ", err)
		return nil, err
	} else if viaCepResponse.Erro {
		slog.Error(fmt.Sprintf("Error invalid address by CEP: %v", viaCepResponse))
		return nil, ErrCEPNotFound
	}

	return &viaCepResponse, nil
}
