<template>
  <transition name="slide-fade">
    <v-container v-if="bingo">
      <v-row align="center" justify="center">
        <v-col cols="12">
          <v-card v-if="bingo && bingo.BingoID">
            <v-card-title> Bingo {{ bingo.name }}</v-card-title>
            <v-card-subtitle>{{ $route.params["identifier"] }}</v-card-subtitle>
            <v-card-text>
              <b>Tableros: {{ bingo.BoardsSold }}</b
              ><br />
              <b
                >Modo de Juego: {{ bingo.CurrentMode }}
                <GameModeChip :currentMode="bingo.CurrentMode" /></b
              ><br />
              <v-row class="mt-4">
                <v-col cols="4">
                  <v-btn
                    :class="
                      !bingo.Playing
                        ? 'green accent-2'
                        : 'red accent-3 white--text'
                    "
                    @click="reqPlayMode(!bingo.Playing)"
                    small
                    block
                  >
                    <template v-if="!bingo.Playing">
                      <v-icon>mdi-play</v-icon> Jugar
                    </template>
                    <template v-else>
                      <v-icon class="mr-2" small>mdi-stop</v-icon> Terminar
                    </template>
                  </v-btn>
                </v-col>
                <v-col cols="4">
                  <v-btn
                    color="blue accent-2 white--text"
                    @click="reqGenerateBoard"
                    small
                    block
                  >
                    Generar Tablero
                  </v-btn>
                </v-col>
                <v-col cols="4">
                  <v-btn
                    color="blue accent-2 white--text"
                    @click="reqCheckBoard(false)"
                    small
                    block
                  >
                    <v-icon class="mr-2" small>mdi-check-circle</v-icon>
                    Verificar Tablero
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      bingo.playing: {{ bingo.Playing }}
      <transition name="slide-fade">
        <v-card v-if="bingo && bingo.Playing">
          <v-card-title> Jugando {{ bingo.Name }} </v-card-title>
          <v-card-text>
            <v-row align="center" justify="center" class="my-4 text-center">
              <v-col cols="6">
                <v-text-field
                  label="Balota"
                  v-model="drawBalot"
                  outlined
                  dense
                  hide-details
                ></v-text-field>
              </v-col>
              <v-col cols="1" class="text-left">
                <v-btn small fab color="blue accent-3" @click="reqDrawBalot">
                  <v-icon color="white"> mdi-ballot </v-icon>
                </v-btn>
              </v-col>
            </v-row>
            <v-row justify="center">
              <v-col cols="4">
                <v-select
                  label="game_mode"
                  solo
                  flat
                  v-model="newGameMode"
                  :items="gameModes"
                  outlined
                  filled
                  dense
                  append-icon="mdi-chevron-down"
                >
                  <template v-slot:selection="{ item }">
                    <GameModeChip :currentMode="item" />
                  </template>
                </v-select>
              </v-col>
            </v-row>
            <v-row justify="center" class="text-center">
              <v-col cols="12" class="title text-center">
                Balotas Jugadas
              </v-col>
              <v-col cols="2">
                <p class="display-1">B</p>
                <template
                  v-for="balot in drawnBalots.map((v) =>
                    v[0] == 'B' ? v : null
                  )"
                >
                  <p v-if="balot" :key="'balotB' + balot">
                    {{ balot }}
                  </p>
                </template>
              </v-col>
              <v-col cols="2">
                <p class="display-1">I</p>
                <template
                  v-for="balot in drawnBalots.map((v) =>
                    v[0] == 'I' ? v : null
                  )"
                >
                  <p v-if="balot" :key="'balotI' + balot">
                    {{ balot }}
                  </p>
                </template>
              </v-col>
              <v-col cols="2">
                <p class="display-1">N</p>
                <template
                  v-for="balot in drawnBalots.map((v) =>
                    v[0] == 'N' ? v : null
                  )"
                >
                  <p v-if="balot" :key="'balotN' + balot">
                    {{ balot }}
                  </p>
                </template>
              </v-col>
              <v-col cols="2">
                <p class="display-1">G</p>
                <template
                  v-for="balot in drawnBalots.map((v) =>
                    v[0] == 'G' ? v : null
                  )"
                >
                  <p v-if="balot" :key="'balotG' + balot">
                    {{ balot }}
                  </p>
                </template>
              </v-col>
              <v-col cols="2">
                <p class="display-1">O</p>
                <template
                  v-for="balot in drawnBalots.map((v) =>
                    v[0] == 'O' ? v : null
                  )"
                >
                  <p v-if="balot" :key="'balotO' + balot">
                    {{ balot }}
                  </p>
                </template>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </transition>
      <v-dialog v-model="dialog_check" max-width="420">
        <v-card>
          <v-card-title class="text-h5"> Verificar tablero </v-card-title>
          <v-card-subtitle>
            Modo de Juego: <GameModeChip :currentMode="bingo.CurrentMode" />
          </v-card-subtitle>
          <v-card-text class="text-center">
            <v-row v-if="!dialog_check_mode" align="center">
              <v-col cols="8">
                <v-text-field
                  label="Tablero"
                  v-model="checkBoardId"
                  outlined
                  dense
                  hide-details
                ></v-text-field>
              </v-col>
              <v-col cols="4">
                <v-btn
                  small
                  color="primary"
                  @click="reqCheckBoard(checkBoardId)"
                  >Verificar</v-btn
                >
              </v-col>
            </v-row>

            <v-circular-progress
              v-if="dialog_check_mode == 'loading'"
              intermitent
            >
            </v-circular-progress>
          </v-card-text>
        </v-card>
      </v-dialog>
    </v-container>
  </transition>
</template>

<script>
import local_api from "~/api/local_api";
import GameModeChip from "~/components/GameModeChip";

export default {
  components: {
    GameModeChip,
  },
  data: () => ({
    gameModes: [
      "lb",
      "li",
      "ln",
      "lg",
      "lo",
      "a",
      "c",
      "o",
      "n",
      "/",
      "\\",
      "l1",
      "l2",
      "l3",
      "l4",
      "l5",
    ],
    newGameMode: "",
    drawBalot: "",
    checkBoardId: "",
    bingo: false,
    generated: false,
    dialog_check: false,
    dialog_check_mode: false,
  }),
  computed: {
    drawnBalots() {
      if (!this.bingo) return [];
      else return this.bingo.DrawnBalots.split(",");
    },
  },
  watch: {
    newGameMode(ngm) {
      console.log("new game mode changes", ngm);
      try {
        local_api
          .api_post("game/setmode", {
            BingoID: this.$route.params["identifier"],
            Param: ngm,
          })
          .then((respuesta) => {
            console.log(respuesta);
            if (respuesta.Status == "OK") {
              this.bingo.CurrentMode = respuesta.Message;
              let ogBingo = this.$deepCopyObject(this.bingo);
              this.bingo = false;
              setTimeout(() => {
                this.bingo = this.$deepCopyObject(ogBingo);
                console.log("now bingo is", this.bingo);
              }, 50);
            }
          });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
  },
  created() {
    console.log(this.$route);
  },
  mounted() {
    console.log(this.$route);
    this.loadGame();
  },
  methods: {
    loadGame() {
      console.log("loadGame", this.$route.params["identifier"]);
      try {
        local_api
          .api_get("game", this.$route.params["identifier"])
          .then((respuesta) => {
            console.log(respuesta);
            this.bingo = respuesta;
            this.newGameMode = this.bingo.CurrentMode;
            // this.bingo.Playing = this.bingo.Playing == 'true'
          });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
    reqGenerateBoard() {
      console.log("reqGenerateBoard");
      try {
        local_api
          .api_post("game/generate", {
            BingoID: this.$route.params["identifier"],
          })
          .then((respuesta) => {
            console.log(respuesta);
            this.generated = respuesta;
          });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
    reqPlayMode(play) {
      console.log("reqPlay", play);
      try {
        local_api
          .api_post("game/setmode", {
            BingoID: this.$route.params["identifier"],
            Param: play ? "PLAY" : "STOP",
          })
          .then((respuesta) => {
            console.log(respuesta);
            if (respuesta.Status == "OK") {
              this.bingo.Playing = play;
              let ogBingo = this.$deepCopyObject(this.bingo);
              this.bingo = false;
              setTimeout(() => {
                this.bingo = this.$deepCopyObject(ogBingo);
                console.log("now bingo is", this.bingo);
              }, 50);
            }
          });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
    reqDrawBalot() {
      let balot = this.drawBalot;
      console.log("reqDrawBalot", balot);
      try {
        local_api
          .api_post("game/drawbalot", {
            BingoID: this.$route.params["identifier"],
            Balot: balot,
          })
          .then((respuesta) => {
            console.log(respuesta);
            if (respuesta.Status == "OK") {
              console.log("appending", respuesta.Message);
              if (this.bingo.DrawnBalots == "") {
                this.bingo.DrawnBalots = respuesta.Message;
              } else {
                this.bingo.DrawnBalots += "," + respuesta.Message;
              }

              console.log(this.drawnBalots);
            }
          });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
    reqCheckBoard(boardID = false) {
      console.log("reqCheckBoard", boardID);
      if (!boardID) {
        this.dialog_check = true;
        return;
      }
      this.dialog_check_mode == "loading";
      try {
        local_api
          .api_post("game/check", {
            BingoID: this.$route.params["identifier"],
            BoardID: boardID,
          })
          .then((respuesta) => {
            console.log(respuesta);
            if (respuesta.Status == "OK") {
              console.log("appending", respuesta.Message);
              if (this.bingo.DrawnBalots == "") {
                this.bingo.DrawnBalots = respuesta.Message;
              } else {
                this.bingo.DrawnBalots += "," + respuesta.Message;
              }

              console.log(this.drawnBalots);
            }
          });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
  },
};
</script>

<style>
/* Enter and leave animations can use different */
/* durations and timing functions.              */
.slide-fade-enter-active {
  transition: all 1s ease;
}
.slide-fade-leave-active {
  transition: all 1s cubic-bezier(1, 0.5, 0.8, 1);
}
.slide-fade-enter, .slide-fade-leave-to
/* .slide-fade-leave-active below version 2.1.8 */ {
  transform: translateX(10px);
  opacity: 0;
}
</style>