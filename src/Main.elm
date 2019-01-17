module Main exposing (ColumnInfo, CraigslistHTML, Model, Msg(..), Url, categorySelector, citySelector, deleteColumnButton, init, loadRefreshButton, main, postBody, queryColumn, queryDecoder, queryGridColumnWrap, queryResults, subscriptions, update, view)

import Browser
import Css exposing (..)
import Html
import Html.Styled exposing (..)
import Html.Styled.Attributes exposing (..)
import Html.Styled.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, field, list, string)
import Json.Encode exposing (..)



-- MAIN


main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view >> toUnstyled
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
    , urlSetId : Int
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
        0
        "dummy debug"
    , httpGETUrlSet "0"
    )



-- UPDATE


type Msg
    = FormInput FormInputElement ColumnId String
    | LoadButtonPressed ColumnId
    | ReceivedQueryResults (Result Http.Error String) ColumnId
    | ReceivedUrlSet (Result Http.Error (List String))
    | IncrementUrlSetNumber
    | DecrementUrlSetNumber
    | AddColumnButtonClicked
    | DeleteButtonPressed ColumnId


type FormInputElement
    = FormUrlInput
    | FormQueryInput
    | FormCategoryInput
    | FormCityInput


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FormInput formInputElement columnId input ->
            let
                updateFunc =
                    case formInputElement of
                        FormUrlInput ->
                            updateColumnInfosFormUrl

                        FormQueryInput ->
                            updateColumnInfosFormQuery

                        FormCategoryInput ->
                            updateColumnInfosFormCategory

                        FormCityInput ->
                            --not implemented
                            updateColumnInfosFormCategory
            in
            ( { model
                | debugBreadcrumb = input
                , columnInfos = updateFunc model.columnInfos columnId input
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
                            , ( "setIndex", Json.Encode.int model.urlSetId )
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
                    ( { model
                        | columnInfos =
                            updateColumnInfosHtml model.columnInfos columnId <| "fail"
                      }
                    , Cmd.none
                    )

        ReceivedUrlSet result ->
            case result of
                Ok urlSet ->
                    let
                        newmodel =
                            { model
                                | debugBreadcrumb = String.concat urlSet
                                , columnInfos = updateColumnInfosNewUrlSet urlSet
                            }
                    in
                    ( newmodel
                    , let
                        f index url =
                            httpRequestColumn url index newmodel.urlSetId
                      in
                      Cmd.batch <| List.indexedMap f urlSet
                      --Cmd.none --To disable loading all URLs after a refresh
                    )

                Err e ->
                    ( { model | debugBreadcrumb = "watfail" }
                    , Cmd.none
                    )

        DecrementUrlSetNumber ->
            let
                newUrlSetId =
                    if model.urlSetId == 0 then
                        0

                    else
                        model.urlSetId - 1
            in
            ( { model
                | urlSetId = newUrlSetId
                , debugBreadcrumb = "the model is " ++ String.fromInt newUrlSetId
              }
            , httpGETUrlSet <| String.fromInt newUrlSetId
            )

        IncrementUrlSetNumber ->
            let
                newUrlSetId =
                    model.urlSetId + 1
            in
            ( { model
                | urlSetId = newUrlSetId
                , debugBreadcrumb = "the model is " ++ String.fromInt newUrlSetId
              }
            , httpGETUrlSet <| String.fromInt newUrlSetId
            )

        AddColumnButtonClicked ->
            ( model
            , httpJSONBodyReceivedUrlSet "PUT" <|
                Json.Encode.object
                    [ ( "setIndex", Json.Encode.int model.urlSetId )
                    ]
            )

        DeleteButtonPressed columnId ->
            ( model
            , httpJSONBodyReceivedUrlSet "DELETE" <|
                Json.Encode.object
                    [ ( "columnIndex", Json.Encode.int columnId )
                    , ( "setIndex", Json.Encode.int model.urlSetId )
                    ]
            )


httpGETUrlSet : String -> Cmd Msg
httpGETUrlSet columnId =
    Http.request
        { method = "GET"
        , url = "http://localhost:8080/api/" ++ columnId
        , body = Http.emptyBody
        , expect = Http.expectJson ReceivedUrlSet getUrlSetDecoder
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        }


httpJSONBodyReceivedUrlSet theMethod json =
    Http.request
        { method = theMethod
        , url = "http://localhost:8080/api/"
        , body =
            Http.jsonBody json
        , expect = Http.expectJson ReceivedUrlSet getUrlSetDecoder
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        }


httpRequestColumn : String -> Int -> Int -> Cmd Msg
httpRequestColumn url columnId setId =
    Http.request
        { method = "POST"
        , body =
            Http.jsonBody <|
                Json.Encode.object
                    [ ( "searchURL", Json.Encode.string url )
                    , ( "columnIndex", Json.Encode.int columnId )
                    , ( "setIndex", Json.Encode.int setId )
                    ]
        , url = "http://localhost:8080/api/"
        , expect = Http.expectJson (\jsonResult -> ReceivedQueryResults jsonResult columnId) queryDecoder
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        }


updateColumnInfoFieldById : List ColumnInfo -> Int -> (ColumnInfo -> ColumnInfo) -> String -> List ColumnInfo
updateColumnInfoFieldById origColumnInfos columnId updateFunc query =
    let
        f g columnInfo =
            if columnInfo.id == columnId then
                g columnInfo

            else
                columnInfo
    in
    List.map (f updateFunc) origColumnInfos


updateColumnInfosHtml : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosHtml origColumnInfos columnId html =
    let
        z columnInfo =
            { columnInfo | responseHtml = html }
    in
    updateColumnInfoFieldById origColumnInfos columnId z html


updateColumnInfosFormQuery : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosFormQuery origColumnInfos columnId query =
    let
        z columnInfo =
            { columnInfo | url = query, formQuery = query }
    in
    updateColumnInfoFieldById origColumnInfos columnId z query


updateColumnInfosFormCategory : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosFormCategory origColumnInfos columnId category =
    let
        z columnInfo =
            { columnInfo | url = columnInfo.formQuery ++ category }
    in
    updateColumnInfoFieldById origColumnInfos columnId z category


updateColumnInfosFormUrl : List ColumnInfo -> Int -> String -> List ColumnInfo
updateColumnInfosFormUrl origColumnInfos columnId urlArg =
    let
        z columnInfo =
            { columnInfo | url = urlArg }
    in
    updateColumnInfoFieldById origColumnInfos columnId z urlArg


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
        [ topHeader model.urlSetId

        --, text model.debugBreadcrumb
        , topTable [] <| List.map queryGridColumnWrap model.columnInfos
        ]


topHeader : Int -> Html Msg
topHeader urlSetId =
    styled div
        [ displayFlex
        , Css.height (vh 5)
        ]
        []
        [ styled h1 [ margin (px 20) ] [] [ text "Craigslist Side-by-Side" ]
        , urlSetView <| String.fromInt urlSetId
        , button [ onClick AddColumnButtonClicked ] [ text "Add Column" ]
        ]


urlSetView : String -> Html Msg
urlSetView setNumber =
    styled div
        []
        []
        [ button [ onClick DecrementUrlSetNumber ] [ text "-" ]
        , styled input [ textAlign center ] [ placeholder "URL SET", value setNumber ] []
        , button [ onClick IncrementUrlSetNumber ] [ text "+" ]
        ]


topTable : List (Attribute msg) -> List (Html msg) -> Html msg
topTable attrs children =
    styled div
        [ overflowX scroll
        , Css.height (vh 95)
        ]
        []
        [ styled Html.Styled.table
            []
            attrs
            [ styled Html.Styled.tr
                []
                []
                children
            ]
        ]


queryGridColumnWrap : ColumnInfo -> Html Msg
queryGridColumnWrap columnInfo =
    styled Html.Styled.td
        []
        []
        [ queryColumn columnInfo ]


queryColumn : ColumnInfo -> Html Msg
queryColumn columnInfo =
    styled div
        [ Css.width auto
        , padding (px 5)
        , overflowY scroll
        , overflowX Css.hidden
        , Css.height (vh 90)
        ]
        []
        [ styled input
            [ display block, Css.width (pct 100) ]
            [ placeholder "URL", value columnInfo.url, onInput (\input -> FormInput FormUrlInput columnInfo.id input) ]
            []
        , styled input
            [ display block, margin (px 10), Css.width (pct 50) ]
            [ placeholder "Search Query", onInput (\input -> FormInput FormQueryInput columnInfo.id input) ]
            []
        , styled div [ display inline, margin (px 10) ] [] [ categorySelector columnInfo.id (\input -> FormInput FormCategoryInput columnInfo.id input) ]
        , styled div [ display inline, margin (px 10) ] [] [ citySelector ]
        , styled div
            [ displayFlex, flexDirection row, padding (px 15), justifyContent spaceBetween ]
            []
            [ styled div [] [] [ loadRefreshButton columnInfo.id ]
            , styled div [] [] [ deleteColumnButton columnInfo.id ]
            ]
        , styled div [ display block ] [] []
        , styled div [ display block ] [] [ queryResults columnInfo.responseHtml ]
        ]


queryResults : String -> Html Msg
queryResults result =
    postBody result


categorySelector : ColumnId -> (String -> Msg) -> Html Msg
categorySelector id callback =
    select [ onInput callback ]
        [ option [] [ text "Select Category" ]
        , option [] [ text "option 2" ]
        ]


citySelector : Html Msg
citySelector =
    select []
        [ option [] [ text "Select City" ]
        , option [] [ text "Birminham" ]
        ]


loadRefreshButton : ColumnId -> Html Msg
loadRefreshButton param =
    button
        [ onClick (LoadButtonPressed param)
        ]
        [ text "Load Results and Save URL" ]


deleteColumnButton : ColumnId -> Html Msg
deleteColumnButton param =
    button
        [ onClick (DeleteButtonPressed param)
        ]
        [ text "Delete Column" ]



-- This rendered-html node is a custom element
-- defined in the html in a <script> tag
-- https://leveljournal.com/server-rendered-html-in-elm


postBody : String -> Html msg
postBody html =
    Html.Styled.node "rendered-html"
        [ Html.Styled.Attributes.property "content" (Json.Encode.string html) ]
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
