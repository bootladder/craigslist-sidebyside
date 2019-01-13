module Main exposing (ColumnInfo, CraigslistHTML, Model, Msg(..), Url, categorySelector, citySelector, deleteColumnButton, init, loadRefreshButton, main, postBody, queryColumn, queryDecoder, queryGridColumnWrap, queryResults, subscriptions, update , view)

import Bootstrap.Button as Button
import Bootstrap.CDN as CDN
import Bootstrap.Form.Input as Input
import Bootstrap.Grid as Grid
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode exposing (Decoder, field, string)
import Json.Encode exposing (..)



-- MAIN


main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        }



-- MODEL


type alias Url =
    String


type alias CraigslistHTML =
    String


type alias ColumnInfo =
    { id : String
    , url : String
    , responseHtml : String
    }


type alias Model =
    { columnInfos : List ColumnInfo
    , debugBreadcrumb : String
    }


init : () -> ( Model, Cmd Msg )
init _ =
    -- The initial model comes from a Request, now it is hard coded
    ( Model [ { id = "1", url = "hardUrl1", responseHtml = "result1" }, { id = "2", url = "hardUrl2", responseHtml = "result2" } ] "dummy debug"
    , Cmd.none
    )



-- UPDATE


type Msg
    = SearchQueryInput String String
    | LoadButtonPressed String
    | ReceivedQueryResults (Result Http.Error String) String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        SearchQueryInput columnId input ->
            ( { model | debugBreadcrumb = input }
            , Cmd.none
            )

        LoadButtonPressed columnId ->
            ( model
            , Http.request
                { method = "POST"
                , body =
                    Http.jsonBody <|
                        Json.Encode.object
                            [ ( "searchURL", Json.Encode.string columnId )
                            , ( "columnIndex", Json.Encode.int 0 )
                            , ( "setIndex", Json.Encode.int 0 )
                            ]
                , url = "http://localhost:8080/api/"
                , expect = Http.expectJson (\result -> ReceivedQueryResults result columnId) queryDecoder
                , headers = []
                , timeout = Nothing
                , tracker = Nothing
                }
            )

        ReceivedQueryResults result columnId ->
            case result of
                Ok fullText ->
                    ( { model
                        | columnInfos =
                            updateColumnInfosHtml model.columnInfos columnId fullText
                      }
                    , Cmd.none
                    )

                Err e ->
                    case e of
                        Http.BadBody s ->
                            ( { model
                                | columnInfos =
                                    updateColumnInfosHtml model.columnInfos columnId <| "fail" ++ s
                              }
                            , Cmd.none
                            )

                        Http.BadUrl _ ->
                            ( model, Cmd.none )

                        Http.Timeout ->
                            ( model, Cmd.none )

                        Http.NetworkError ->
                            ( model, Cmd.none )

                        Http.BadStatus _ ->
                            ( model, Cmd.none )


updateColumnInfosHtml : List ( ColumnInfo ) -> String -> String -> List ( ColumnInfo )
updateColumnInfosHtml origColumnInfos columnId html =
    let
        f columnInfo =
            if columnInfo.id == columnId then
                {
                    id =  columnInfo.id,
                    url = columnInfo.url,
                    responseHtml = html
                }

            else
                columnInfo
    in
    List.map f origColumnInfos



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ Grid.container []
            [ CDN.stylesheet
            , text model.debugBreadcrumb
            , Grid.row [] <| List.map queryGridColumnWrap model.columnInfos
            ]
        ]


queryGridColumnWrap columnInfo =
    Grid.col [] [ queryColumn columnInfo ]


queryColumn : ColumnInfo -> Html Msg
queryColumn columnInfo =
    Grid.container []
        [ Grid.row []
            [ Grid.col []
                [ Input.text [ Input.attrs [ placeholder "URL" ] ] ]
            ]
        , Grid.row []
            [ Grid.col []
                [ Input.text [ Input.attrs [ placeholder "Search Query", onInput (SearchQueryInput columnInfo.id ) ] ] ]
            ]
        , Grid.row [] [ Grid.col [] [ categorySelector ] ]
        , Grid.row [] [ Grid.col [] [ citySelector ] ]
        , Grid.row []
            [ Grid.col []
                [ loadRefreshButton columnInfo.id
                , deleteColumnButton columnInfo.id
                ]
            ]
        , Grid.row []
            [ Grid.col []
                [ queryResults columnInfo.responseHtml
                ]
            ]
        ]


queryResults : String -> Html Msg
queryResults result =
    postBody result


categorySelector : Html Msg
categorySelector =
    select []
        [ option [] [ text "Select Category" ]
        , option [] [ text "option 2" ]
        ]


citySelector : Html Msg
citySelector =
    select []
        [ option [] [ text "Select City" ]
        , option [] [ text "Birminham" ]
        ]


loadRefreshButton : String -> Html Msg
loadRefreshButton param =
    Button.button
        [ Button.primary
        , Button.small
        , Button.block
        , Button.onClick (LoadButtonPressed param)
        ]
        [ text "Load Results and Save URL" ]


deleteColumnButton : String -> Html Msg
deleteColumnButton param =
    Button.button
        [ Button.danger
        , Button.small
        , Button.block
        , Button.onClick (LoadButtonPressed param)
        ]
        [ text "Delete this column" ]



-- This rendered-html node is a custom element
-- defined in the html in a <script> tag
-- https://leveljournal.com/server-rendered-html-in-elm


postBody : String -> Html msg
postBody html =
    Html.node "rendered-html"
        [ property "content" (Json.Encode.string html) ]
        []



-- HTTP


queryDecoder : Decoder String
queryDecoder =
    field "response" Json.Decode.string
