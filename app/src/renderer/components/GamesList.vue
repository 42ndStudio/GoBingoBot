<template>
  <v-card>
    <v-card-title>
      <v-toolbar color="indigo" dark>
        <v-toolbar-title>Mis Bingos</v-toolbar-title>

        <v-spacer></v-spacer>

        <v-btn icon @click="dialog_add = true">
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </v-toolbar>
    </v-card-title>
    <v-card-text>
      <v-list two-line>
        <template v-for="game in games">
          <v-list-item :key="'item' + game.name" :to="'/bingo/' + game.id">
            <v-list-item-content>
              <v-list-item-title>{{ game.name }}</v-list-item-title>
              <v-list-item-subtitle
                >{{ game.boards }} tableros.</v-list-item-subtitle
              >
            </v-list-item-content>
          </v-list-item>
        </template>
      </v-list>
    </v-card-text>
    <v-dialog v-model="dialog_add" max-width="420">
      <v-card>
        <v-card-title class="text-h5"> Organizar nuevo bingo </v-card-title>

        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12">
                <v-text-field
                  label="Nombre"
                  hint="dale un nombre para identificar este juego."
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-text-field
                  label="Cartones"
                  type="number"
                  hint="también los puedes generar más tarde."
                ></v-text-field>
              </v-col>
              <v-col cols="12">
                <v-file-input
                  accept="image/*"
                  label="Plantilla"
                  hint="puedes configurarlo más adelante"
                ></v-file-input>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <v-btn color="green darken-1" text @click="dialog_add = false">
            Cancelar
          </v-btn>

          <v-btn color="green darken-1" text @click="submitNewGame">
            Crear
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
import { dialog } from 'electron';
export default {
  data: () => ({
    dialog_add: false,
    games: [
      {
        id: "bingo_1",
        name: "Bingo 1",
        boards: 420,
        status: "planned",
      },
      {
        id: "bingo_2",
        name: "Mi Bingo",
        boards: 42,
        status: "played",
      },
    ],
  }),
  methods: {
    submitNewGame() {
      console.log("saving game");
      this.dialog_add = false;
    },
  },
};
</script>

<style>
</style>