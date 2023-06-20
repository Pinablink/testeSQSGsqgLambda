package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testeSQSGsqgLambda/cadastro"

	"github.com/Pinablink/sqg"
	"github.com/Pinablink/sqg/util"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//
func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var nameQueue string = os.Getenv("NAME_QUEUE")
	var refSqg *sqg.SQG = sqg.NewSQG(nameQueue)

	if request.HTTPMethod == "POST" {
		return postInQueue(request, refSqg, nameQueue)
	} else if request.HTTPMethod == "GET" {
		return getInQueue(request, refSqg, nameQueue)
	} else {
		messageresponse := "Verbo HTTP inválido para essa chamada"
		responseObj := cadastro.CadastroStatus{
			Status:  "ERROR",
			Message: messageresponse,
		}

		apiResponse := facResponse(responseObj, nil)

		return apiResponse, nil
	}

}

// getInQueue: Obtêm a mensagem disponível na fila
func getInQueue(request events.APIGatewayProxyRequest, sqgref *sqg.SQG, nameQueue string) (events.APIGatewayProxyResponse, error) {
	var apiResponse events.APIGatewayProxyResponse
	var responseObj interface{}
	var formatMessageResponse string = "%s. %s"
	var strMessageResponse string
	var refCadastro *cadastro.Cadastro = &cadastro.Cadastro{}
	var refCadastroMsgHeader *cadastro.CadastroMsgHeader = &cadastro.CadastroMsgHeader{}

	reviewMessageOK, deleteMessageOK, err := sqgref.GetMsgInQueue(refCadastro, refCadastroMsgHeader)

	if reviewMessageOK {

		if err != nil {
			strMessageResponse = "Ocorreu um erro na obtenção de dados na fila."
			strMessageResponse = fmt.Sprintf(formatMessageResponse, strMessageResponse, err.Error())
		} else if !reviewMessageOK {
			strMessageResponse = "Não foi encontrado mensagem na fila."
			strMessageResponse = fmt.Sprintf(formatMessageResponse, strMessageResponse, "")
		} else if reviewMessageOK && !deleteMessageOK {
			strMessageResponse = "Mensagem na fila não foi deletada"
			strMessageResponse = fmt.Sprintf(formatMessageResponse, strMessageResponse, "")
		} else if reviewMessageOK {
			strMessageResponse = "Mensagem obtida"
			strMessageResponse = fmt.Sprintf(formatMessageResponse, strMessageResponse, "")
		}

		responseObj = cadastro.TesteCadastro{
			CadastroHeader: *refCadastroMsgHeader,
			DataCadastro:   *refCadastro,
		}

	} else {

		responseObj = cadastro.TesteCadastroResponseMessage{
			Message: "Não existe mensagem na fila para consulta",
		}

	}

	apiResponse = facResponse(responseObj, nil)

	return apiResponse, nil
}

// postInQueue: Adiciona Mensagem a uma fila SQS AWS
// request - Struct de requisição do API Gateway
// refSqg - Struct com os dados da mensagem
// nameQueue - Nome da fila
func postInQueue(request events.APIGatewayProxyRequest, refSqg *sqg.SQG, nameQueue string) (events.APIGatewayProxyResponse, error) {
	var apiResponse events.APIGatewayProxyResponse
	var responseObj cadastro.CadastroStatus

	err, cadastroRequest := getRequestData(request)

	if err != nil {

		messageresponse := "Erro ao incluir cadastro"

		responseObj = cadastro.CadastroStatus{
			Status:  "ERROR",
			Message: messageresponse,
		}

		err = errors.New(messageresponse)

	} else {

		if len(nameQueue) == 0 {

			messageresponse := "Identificacao da fila não encontrado"
			responseObj = cadastro.CadastroStatus{
				Status:  "ERROR",
				Message: messageresponse,
			}

			err = errors.New(messageresponse)

		} else {

			headerMessage := cadastroRequest.CadastroHeader
			mcadastro := cadastroRequest.DataCadastro

			var messageModel util.GSQGMessageModel = util.GSQGMessageModel{
				ContentMessage: mcadastro,
				DataMessage:    headerMessage,
			}

			refSqg.SetGSQGMessageModel(messageModel)

			responseId, responseErr := refSqg.JoinTheQueue()

			if responseErr != nil {
				messageError := "Erro na inclusao de dados na fila\n %s"
				messageError = fmt.Sprintf(messageError, responseErr)

				responseObj = cadastro.CadastroStatus{
					Status:  "ERROR",
					Message: messageError,
				}

				err = errors.New(messageError)

			} else {

				responseObj = cadastro.CadastroStatus{
					Status:     "OK",
					Message:    "Sucesso ao incluir cadastro na fila",
					IdResponse: *responseId,
				}

			}

		}

	}

	apiResponse = facResponse(responseObj, err)

	return apiResponse, err
}

//
func getRequestData(request events.APIGatewayProxyRequest) (error, cadastro.TesteCadastro) {
	byteBody := []byte(request.Body)
	var cadastroRequest cadastro.TesteCadastro
	err := json.Unmarshal(byteBody, &cadastroRequest)
	return err, cadastroRequest
}

//
func facResponse(structResponse interface{}, refError error) events.APIGatewayProxyResponse {
	var mapHeader map[string]string
	var dataByte []byte
	var dataReturn string
	var statusCode int

	statusCode = 200
	mapHeader = make(map[string]string)
	mapHeader["Content-Type"] = "application/json"

	dataByte, refError = json.Marshal(structResponse)

	if refError != nil {
		statusCode = 500
	}

	dataReturn = string(dataByte)

	return events.APIGatewayProxyResponse{
		Body:       dataReturn,
		StatusCode: statusCode,
		Headers:    mapHeader,
	}
}

func main() {
	lambda.Start(HandleRequest)
}
