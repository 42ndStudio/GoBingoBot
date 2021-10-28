<template>
  <v-container>
    <v-row align="center" justify="center">
      <v-col cols="12">
        <v-card>
          <v-card-title> Bingo {{ $route.params['identifier'] }} </v-card-title>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  mounted() {
    console.log(this.$route);
  },
};
</script>

<style>
</style>