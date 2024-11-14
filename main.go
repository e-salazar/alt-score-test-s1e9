package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var (
	m_liquid float64 //Pendiente
	b_liquid float64 //Corte en Y

	m_vapor float64 //Pendiente
	b_vapor float64 //Corte en Y
)

type PhaseChangeDiagramResponse struct {
	SpecificVolumeLiquid float64 `json:"specific_volume_liquid"`
	SpecificVolumeVapor  float64 `json:"specific_volume_vapor"`
}

func main() {
	liquidP1x := 0.00105
	liquidP1y := 0.05
	liquidP2x := 0.0035
	liquidP2y := 10.0

	vaporP1x := 30.0
	vaporP1y := 0.05
	vaporP2x := 0.0035
	vaporP2y := 10.0

	m_liquid = (liquidP2y - liquidP1y) / (liquidP2x - liquidP1x)
	b_liquid = liquidP1y - (m_liquid * liquidP1x)

	m_vapor = (vaporP2y - vaporP1y) / (vaporP2x - vaporP1x)
	b_vapor = vaporP1y - (m_vapor * vaporP1x)

	http.HandleFunc("/phase-change-diagram", phaseChangeDiagramHandler)

	fmt.Println("API ejecutándose en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Endpoint /phase-change-diagram
func phaseChangeDiagramHandler(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	if httpRequest.Method != http.MethodGet {
		http.Error(httpResponse, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	pressureParam := httpRequest.URL.Query().Get("pressure")
	if pressureParam == "" {
		http.Error(httpResponse, "El parámetro 'pressure' es requerido", http.StatusBadRequest)
		return
	}

	pressure, err := strconv.ParseFloat(pressureParam, 64)
	if err != nil {
		http.Error(httpResponse, "Parámetro 'pressure' debe ser un número", http.StatusBadRequest)
		return
	}

	liquidVolume, err := calculateLiquidVolume(pressure)
	if err != nil {
		http.Error(httpResponse, "Error en el cálculo de liquidVolume", http.StatusInternalServerError)
		return
	}

	vaporVolume, err := calculateVaporVolume(pressure)
	if err != nil {
		http.Error(httpResponse, "Error en el cálculo de vaporVolume", http.StatusInternalServerError)
		return
	}

	response := PhaseChangeDiagramResponse{
		SpecificVolumeLiquid: liquidVolume,
		SpecificVolumeVapor:  vaporVolume,
	}
	httpResponse.Header().Set("Content-Type", "application/json")
	json.NewEncoder(httpResponse).Encode(response)
}

// Ecuación de la recta: y = m * x + b
// =>	pressure = m * volume + b
// =>	volume = (pressure - b) / m

func calculateLiquidVolume(pressure float64) (float64, error) {
	if m_liquid == 0 {
		return 0, fmt.Errorf("pendiente m_liquid no puede ser cero")
	}
	return (pressure - b_liquid) / m_liquid, nil
}

func calculateVaporVolume(pressure float64) (float64, error) {
	if m_vapor == 0 {
		return 0, fmt.Errorf("pendiente m_vapor no puede ser cero")
	}
	return (pressure - b_vapor) / m_vapor, nil
}
