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
    { name : String
    , password : String
    , passwordAgain : String
    , queryResult : String
    }

init : () -> ( Model, Cmd Msg)
init _ =
    (Model "" "" "" ""
    , Cmd.none
    )

-- UPDATE
type Msg
  = Name String
  | Password String
  | PasswordAgain String
  | LoadButtonPressed String
  | ReceivedQueryResults (Result Http.Error String)

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    Name name ->
      ({ model | name = name }
      , Cmd.none)

    Password password ->
      ({ model | password = password }
      , Cmd.none)

    PasswordAgain password ->
      ({ model | passwordAgain = password }
      , Cmd.none)

    LoadButtonPressed columnId ->
      ({ model | queryResult = "helloresult" ++ columnId }
      , Http.request
      {
        method = "POST"
      , body = Http.jsonBody <|
              Json.Encode.object [
                   ( "searchURL", Json.Encode.string "mytitle")
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
        [ CDN.stylesheet
          ,Grid.row [] <| List.repeat 5 (Grid.col [] [queryColumn model])
        ]
    , validator model
    ]


queryColumn: Model -> Html Msg
queryColumn model =
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
                   loadRefreshButton "1234567890"
                  ,deleteColumnButton "1234567890"
                 ]
            ]
        , Grid.row [] [Grid.col [] [ queryResults model ]]
            ]

queryResults : Model -> Html Msg
queryResults model = postBody model.queryResult

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
        , Button.onClick (PasswordAgain param)
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

 ---------------------------------------------------

 ---------------------------------------------------
validator model =
    div [] [
         viewInput "text" "Name" model.name Name
        , viewInput "password" "Password" model.password Password
        , viewInput "password" "Re-enter Password"
             model.passwordAgain PasswordAgain
        , viewValidation model
        ]

viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ type_ t, placeholder p, value v, onInput toMsg ] []


viewValidation : Model -> Html msg
viewValidation model =
  if model.password == model.passwordAgain then
    div [ style "color" "green" ] [ text "OK" ]
  else
    div [ style "color" "red" ] [ text "Passwords do not match!" ]

