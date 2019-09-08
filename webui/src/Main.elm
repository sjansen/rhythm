module Main exposing (Model(..), Msg(..), init, main, subscriptions, update, view)

import Browser
import Css exposing (..)
import Html.Styled exposing (Html, div, h1, text, toUnstyled)
import Html.Styled.Attributes exposing (css)
import Http



-- MAIN


main =
    Browser.element
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view >> toUnstyled
        }



-- MODEL


type Model
    = Failure
    | Loading
    | Success String


init : () -> ( Model, Cmd Msg )
init _ =
    ( Loading
    , Http.get
        { url = "/api/secret"
        , expect = Http.expectString GotText
        }
    )



-- UPDATE


type Msg
    = GotText (Result Http.Error String)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotText result ->
            case result of
                Ok fullText ->
                    ( Success fullText, Cmd.none )

                Err _ ->
                    ( Failure, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    case model of
        Failure ->
            text "Load failed."

        Loading ->
            text "Loading..."

        Success fullText ->
            div
                [ css
                    [ padding (px 20)
                    ]
                ]
                [ h1 [] [ text "Rhythm" ]
                , div
                    [ css
                        [ backgroundColor (rgb 255 255 255)
                        , border3 (px 1) solid (rgb 120 120 120)
                        , padding (px 10)
                        ]
                    ]
                    [ text fullText ]
                ]
