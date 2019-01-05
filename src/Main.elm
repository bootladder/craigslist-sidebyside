import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onInput)
import Json.Encode exposing (..)
import Bootstrap.CDN as CDN
import Bootstrap.Grid as Grid
import Bootstrap.Button as Button

-- MAIN
main =
  Browser.sandbox { init = init, update = update, view = view }

-- MODEL
type alias Model =
  { name : String
  , password : String
  , passwordAgain : String
  }
init : Model
init =
  Model "" "" ""

-- UPDATE
type Msg
  = Name String
  | Password String
  | PasswordAgain String

update : Msg -> Model -> Model
update msg model =
  case msg of
    Name name ->
      { model | name = name }

    Password password ->
      { model | password = password }

    PasswordAgain password ->
      { model | passwordAgain = password }

-- VIEW
view : Model -> Html Msg
view model =
  div []
    [
      Grid.container []
        [ CDN.stylesheet
          ,Grid.row [] [
                Grid.col [] [ queryColumn ]
               ,Grid.col [] [ queryColumn ]
               ,Grid.col [] [ queryColumn ]
               ,Grid.col [] [ queryColumn ]
              ]
        ]
    , viewInput "text" "Name" model.name Name
    , helloHtml ("yo")
    , text "hello2"
    , viewInput "password" "Password" model.password Password
    , viewInput "password" "Re-enter Password" model.passwordAgain PasswordAgain
    , viewValidation model
    ]

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
                   Button.button [] [text "Load Results and Save URL"]
                  ,Button.button [] [text "Delete this column"]
                 ]
            ]
        , Grid.row [] [Grid.col [] [ queryResults ]]
            ]
queryResults = text "query results"

categorySelector =  select []
                      [ option [] [text "Select Category"]
                      , option [] [text "option 2"]
                      ]

citySelector = select [] [
                 option [] [text "Select City"]
                ,option [] [text "Birminham"]
               ]

helloHtml : String -> Html msg
helloHtml str = text (str)

viewInput : String -> String -> String -> (String -> msg) -> Html msg
viewInput t p v toMsg =
  input [ type_ t, placeholder p, value v, onInput toMsg ] []


viewValidation : Model -> Html msg
viewValidation model =
  if model.password == model.passwordAgain then
    div [ style "color" "green" ] [ text "OK" ]
  else
    div [ style "color" "red" ] [ text "Passwords do not match!" ]
