port module PortExamples exposing (..)

import Html exposing (..)
import Html.Events exposing (..)


type alias Model =
    String


view : Model -> Html Msg
view model =
    div []
        [ button [ onClick SendDataToJS ]
            [ text "Send Data to JavaScript" ]
        , br [] []
        , br [] []
        , text ("Data received from JavaScript: " ++ model)
        ]


type Msg
    = SendDataToJS
    | ReceivedDataFromJS Model


port sendData : String -> Cmd msg


port receiveData : (Model -> msg) -> Sub msg


subscriptions : Model -> Sub Msg
subscriptions _ =
    receiveData ReceivedDataFromJS


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SendDataToJS ->
            ( model, sendData "Hello JavaScript!" )

        ReceivedDataFromJS data ->
            ( data, Cmd.none )


init : ( Model, Cmd Msg )
init =
    ( "", Cmd.none )


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
