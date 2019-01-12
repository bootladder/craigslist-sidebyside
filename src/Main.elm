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
    }

init : () -> ( Model, Cmd Msg)
init _ =
    (Model "" "" ""
    , Cmd.none
    )

-- UPDATE
type Msg
  = Name String
  | Password String
  | PasswordAgain String

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
          ,Grid.row [] (List.repeat 5 (Grid.col [] [queryColumn]))
        ]
    , validator model
    ]


queryColumn: Html Msg
queryColumn =
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
        , Grid.row [] [Grid.col [] [ queryResults ]]
            ]

queryResults : Html Msg
queryResults = text "query results"

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
           , Button.onClick (Password param)
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

