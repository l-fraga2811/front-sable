package supabase

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

// Numeric é um tipo customizado para lidar com campos NUMERIC do PostgreSQL.
//
// PROBLEMA:
// O tipo NUMERIC do PostgreSQL pode ser retornado como string ou number
// dependendo da precisão e do driver. O PostgREST do Supabase às vezes
// retorna como string para preservar precisão.
//
// SOLUÇÃO:
// Criamos um tipo que implementa json.Unmarshaler e aceita ambos os formatos.
//
// CONCEITO: Custom JSON Unmarshaling
// Em Go, podemos customizar como um tipo é deserializado do JSON
// implementando a interface json.Unmarshaler com o método UnmarshalJSON.
type Numeric float64

// UnmarshalJSON implementa json.Unmarshaler para aceitar number ou string.
//
// FLUXO:
//  1. Tenta parsear como json.Number (formato numérico padrão)
//  2. Se falhar, tenta parsear como string
//  3. Se ambos falharem, retorna erro
//
// CONCEITO: Interface em Go
// json.Unmarshaler é uma interface com um único método:
//
//	type Unmarshaler interface {
//	    UnmarshalJSON([]byte) error
//	}
//
// Qualquer tipo que implemente esse método pode customizar sua deserialização.
func (n *Numeric) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*n = 0
		return nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(data, &asNumber); err == nil {
		f, err := asNumber.Float64()
		if err != nil {
			return err
		}
		*n = Numeric(f)
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		f, err := strconv.ParseFloat(asString, 64)
		if err != nil {
			return err
		}
		*n = Numeric(f)
		return nil
	}

	return errors.New("valor numérico inválido")
}
