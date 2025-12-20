// Package models contém as estruturas de dados (structs) da aplicação.
//
// CONCEITO: Structs em Go
// Structs são similares a classes em outras linguagens, mas com diferenças importantes:
//   - Go NÃO tem herança - usa composição (embedding) ao invés
//   - Go NÃO tem construtores - usamos "constructor functions" (ex: NewUser)
//   - Campos públicos começam com MAIÚSCULA, privados com minúscula
//
// CONCEITO: Tags de Struct
// Tags são metadados associados aos campos. Exemplos:
//   - `json:"nome"` - define o nome do campo no JSON
//   - `json:"-"` - ignora o campo na serialização JSON
//   - `json:",omitempty"` - omite se o valor for zero/vazio
//   - `binding:"required"` - validação do Gin (campo obrigatório)
package models

// UserResponse é a struct usada para retornar dados do usuário.
//
// PADRÃO: Separar Request de Response
// Usamos structs diferentes para entrada (Request) e saída (Response) porque:
//   - Request pode ter campos que não queremos retornar (ex: senha)
//   - Response pode ter campos calculados ou formatados
//   - Facilita versionamento da API (mudar response sem afetar request)
//
// NOTA: Com Supabase Auth, os dados do usuário vêm do token JWT.
// Esta struct é usada apenas para formatar a resposta do endpoint /profile.
type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
