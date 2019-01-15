module Main exposing (ColumnInfo, CraigslistHTML, Model, Msg(..), Url, categorySelector, citySelector, deleteColumnButton, init, loadRefreshButton, main, postBody, queryColumn, queryDecoder, queryGridColumnWrap, queryResults, subscriptions, update, view)

--import Bootstrap.Button as Button
--import Bootstrap.CDN as CDN
--import Bootstrap.Form.Input as Input
--import Bootstrap.Form.Select as Select
--import Bootstrap.Grid as Grid
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode exposing (Decoder, field, list, string)
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


type alias ColumnId =
    Int


type alias Url =
    String


type alias CraigslistHTML =
    String


type alias ColumnInfo =
    { id : Int
    , url : String
    , responseHtml : String
    , formQuery : String
    , formCategory : String
    , formCity : String
    }


type alias Model =
    { columnInfos : List ColumnInfo
    , debugBreadcrumb : String
    }



-- INIT


init : () -> ( Model, Cmd Msg )
init _ =
    -- The initial model comes from a Request, now it is hard coded
    ( Model
        [ { id = 0, url = "hardUrl0", responseHtml = "result0", formQuery = "", formCategory = "", formCity = "" }
        , { id = 1, url = "hardUrl1", responseHtml = "result1", formQuery = "", formCategory = "", formCity = "" }
        ]
        "dummy debug"
    , Http.request
        { method = "GET"
        , url = "http://localhost:8080/api/0"
        , body = Http.emptyBody
        , expect = Http.expectJson ReceivedUrlSet getUrlSetDecoder
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        }
    )



-- UPDATE


type Msg
    = UrlInput ColumnId String
    | SearchQueryInput ColumnId String
    | CategoryInput ColumnId String
    | CityInput ColumnId String
    | LoadButtonPressed ColumnId
    | ReceivedQueryResults (Result Http.Error String) ColumnId
    | ReceivedUrlSet (Result Http.Error (List String))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UrlInput columnId input ->
            ( { model
                | debugBreadcrumb = input
                , columnInfos = updateColumnInfosFormUrl model.columnInfos columnId input
              }
            , Cmd.none
            )

        SearchQueryInput columnId input ->
            ( { model
                | debugBreadcrumb = input
                , columnInfos = updateColumnInfosFormQuery model.columnInfos columnId input
              }
            , Cmd.none
            )

        CategoryInput columnId input ->
            ( { model
                | debugBreadcrumb = input
                , columnInfos = updateColumnInfosFormCategory model.columnInfos columnId input
              }
            , Cmd.none
            )

        CityInput columnId input ->
            ( { model
                | debugBreadcrumb = input
                , columnInfos = updateColumnInfosFormQuery model.columnInfos columnId input
              }
            , Cmd.none
            )

        LoadButtonPressed columnId ->
            ( model
            , Http.request
                { method = "POST"
                , body =
                    Http.jsonBody <|
                        Json.Encode.object
                            [ ( "searchURL", Json.Encode.string <| modelGetUrlFromId model columnId )
                            , ( "columnIndex", Json.Encode.int columnId )
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

        ReceivedUrlSet result ->
            case result of
                Ok urlSet ->
                    ( { model
                        | debugBreadcrumb = String.concat urlSet
                        , columnInfos = updateColumnInfosNewUrlSet urlSet
                      }
                    , Cmd.none
                    )

                Err e ->
                    ( { model | debugBreadcrumb = "watfail" }
                    , Cmd.none
                    )


updateColumnInfosHtml : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosHtml origColumnInfos columnId html =
    let
        f columnInfo =
            if columnInfo.id == columnId then
                { id = columnInfo.id
                , url = columnInfo.url
                , responseHtml = html
                , formQuery = columnInfo.formQuery
                , formCategory = columnInfo.formCategory
                , formCity = columnInfo.formCity
                }

            else
                columnInfo
    in
    List.map f origColumnInfos


updateColumnInfosFormQuery : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosFormQuery origColumnInfos columnId query =
    let
        f columnInfo =
            if columnInfo.id == columnId then
                { id = columnInfo.id
                , url = query
                , responseHtml = columnInfo.responseHtml
                , formQuery = query
                , formCategory = columnInfo.formCategory
                , formCity = columnInfo.formCity
                }

            else
                columnInfo
    in
    List.map f origColumnInfos


updateColumnInfosFormCategory : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosFormCategory origColumnInfos columnId category =
    let
        f columnInfo =
            if columnInfo.id == columnId then
                { id = columnInfo.id
                , url = columnInfo.formQuery ++ category
                , responseHtml = columnInfo.responseHtml
                , formQuery = columnInfo.formQuery
                , formCategory = columnInfo.formCategory
                , formCity = columnInfo.formCity
                }

            else
                columnInfo
    in
    List.map f origColumnInfos


updateColumnInfosFormUrl : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosFormUrl origColumnInfos columnId urlArg =
    let
        f columnInfo =
            if columnInfo.id == columnId then
                { id = columnInfo.id
                , url = urlArg
                , responseHtml = columnInfo.responseHtml
                , formQuery = columnInfo.formQuery
                , formCategory = columnInfo.formCategory
                , formCity = columnInfo.formCity
                }

            else
                columnInfo
    in
    List.map f origColumnInfos


updateColumnInfosNewUrlSet : List String -> List ColumnInfo
updateColumnInfosNewUrlSet urls =
    let
        f index str =
            { id = index
            , url = str
            , responseHtml = ""
            , formQuery = ""
            , formCategory = ""
            , formCity = ""
            }
    in
    List.indexedMap f urls


modelGetUrlFromId : Model -> Int -> String
modelGetUrlFromId model columnId =
    let
        l =
            List.filter (\c -> c.id == columnId) model.columnInfos
    in
    case List.head l of
        Just c ->
            c.url

        Nothing ->
            "http://google.com"



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [  
             text model.debugBreadcrumb
            , div [] <| List.map queryGridColumnWrap model.columnInfos
            
        ]


queryGridColumnWrap columnInfo =
     queryColumn columnInfo 


queryColumn : ColumnInfo -> Html Msg
queryColumn columnInfo =
    div [] [
      --Input.text [ Input.attrs [ placeholder "URL", value columnInfo.url, onInput (UrlInput columnInfo.id) ] ] 
            
       --Input.text [ Input.attrs [ placeholder "Search Query", onInput (SearchQueryInput columnInfo.id) ] ] 
            
         categorySelector columnInfo.id 
        , citySelector
        , loadRefreshButton columnInfo.id
                , deleteColumnButton columnInfo.id
                
            
        , queryResults columnInfo.responseHtml
    ]
                
            
        


queryResults : String -> Html Msg
queryResults result =
    postBody result


categorySelector : ColumnId -> Html Msg
categorySelector id =
    text "categorySelectr"
    --Select.select [ Select.attrs [ onInput (CategoryInput id) ] ]
    --    [ Select.item [] [ text "Select Category" ]
    --    , Select.item [] [ text "option 2" ]
    --    ]


citySelector : Html Msg
citySelector =
    select []
        [ option [] [ text "Select City" ]
        , option [] [ text "Birminham" ]
        ]


loadRefreshButton : ColumnId -> Html Msg
loadRefreshButton param =
    text "loadRefreshButton"
    --Button.button
    --    [ Button.primary
    --    , Button.small
    --    , Button.block
    --    , Button.onClick (LoadButtonPressed param)
    --    ]
    --    [ text "Load Results and Save URL" ]


deleteColumnButton : ColumnId -> Html Msg
deleteColumnButton param =
    text "deleteColumnbutton"
  --  Button.button
  --      [ Button.danger
  --      , Button.small
  --      , Button.block
  --      , Button.onClick (LoadButtonPressed param)
  --      ]
  --      [ text "Delete this column" ]



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


getUrlSetDecoder : Decoder (List String)
getUrlSetDecoder =
    field "urls" listStringDecoder


listStringDecoder : Decoder (List String)
listStringDecoder =
    Json.Decode.list Json.Decode.string