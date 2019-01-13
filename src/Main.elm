module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput,onClick)
import Json.Encode exposing (..)
import Bootstrap.CDN as CDN
import Bootstrap.Grid as Grid
import Bootstrap.Button as Button
import Bootstrap.Utilities.Spacing as Spacing
import Http
import Json.Encode
import Json.Decode exposing (Decoder,string,field)

-- MAIN
main =
  Browser.element { init = init, update = update,
                        subscriptions = subscriptions,
                        view = view }

-- MODEL
type alias Model =
    {
     queryResult : String,
     urlResultTuples : (List (String,String))
    }

init : () -> ( Model, Cmd Msg)
init _ =
    (Model "" [("url1","result1"),("url2","result2")]
    , Cmd.none
    )

-- UPDATE
type Msg
  =
    LoadButtonPressed String
  | ReceivedQueryResults (Result Http.Error String)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of

    LoadButtonPressed columnId ->
      ({ model | queryResult = "helloresult" ++ columnId }
      , Http.request
      {
        method = "POST"
      , body = Http.jsonBody <|
              Json.Encode.object [
                   ( "searchURL", Json.Encode.string columnId)
                  ,( "columnIndex", Json.Encode.int 0)
                  ,( "setIndex", Json.Encode.int 0)
                  ]
      , url = "http://localhost:8080/api/"
      , expect = Http.expectJson ReceivedQueryResults queryDecoder
      , headers = []
      , timeout = Nothing
      , tracker = Nothing
      })

    ReceivedQueryResults result ->
      case result of
        Ok fullText ->
          ({ model | queryResult = fullText }, Cmd.none)

        Err e ->
            case e of
                Http.BadBody s ->
                    ({ model | queryResult = "fail"++s }, Cmd.none)

                Http.BadUrl _     ->  (model,Cmd.none)
                Http.Timeout      ->  (model,Cmd.none)
                Http.NetworkError ->  (model,Cmd.none)
                Http.BadStatus _  ->  (model,Cmd.none)


-- SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none

-- VIEW
view : Model -> Html Msg
view model =
  div []
    [
      Grid.container []
        [   CDN.stylesheet
          , Grid.row [] <| List.map queryGridColumnWrap model.urlResultTuples
        ]
    ]


queryGridColumnWrap tuple = Grid.col [] [queryColumn tuple]

queryColumn: (String,String) -> Html Msg
queryColumn urlResultTuple =
    Grid.container []
        [
         Grid.row []
            [ Grid.col []
                [ input [ placeholder "URL" ] [] ]
            ]
        , Grid.row []
            [ Grid.col []
                [ input [ placeholder "Search Query"] [] ]
            ]
        , Grid.row [] [ Grid.col [] [ categorySelector ] ]
        , Grid.row [] [ Grid.col [] [ citySelector ] ]
        , Grid.row []
            [
             Grid.col []
                 [
                   loadRefreshButton  <| Tuple.first(urlResultTuple)
                  ,deleteColumnButton <| Tuple.first(urlResultTuple)
                 ]
            ]
        , Grid.row [] [ Grid.col [] [ 
            queryResults <| Tuple.second(urlResultTuple)
            ]
            ]
        ]

queryResults : String -> Html Msg
queryResults result = postBody result

categorySelector : Html Msg
categorySelector =  select []
                      [ option [] [text "Select Category"]
                      , option [] [text "option 2"]
                      ]

citySelector : Html Msg
citySelector = select [] [
                 option [] [text "Select City"]
                ,option [] [text "Birminham"]
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
    [text "Delete this column"]


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

