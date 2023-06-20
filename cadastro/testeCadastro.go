package cadastro

//
type CadastroMsgHeader struct {
	Titulo      string `json:"titulo" SQGS:"Titulo"`
	Id_Inclusao string `json:"id_inclusao" SQGS:"Id_Inclusao"`
}

//
type Cadastro struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

//
type TesteCadastro struct {
	CadastroHeader CadastroMsgHeader `json:"header"`
	DataCadastro   Cadastro          `json:"cadastro"`
}

//
type TesteCadastroResponseMessage struct {
	Message string `json:"message"`
}

//
type CadastroStatus struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	IdResponse string `json:"idresponse_message_queue"`
}
