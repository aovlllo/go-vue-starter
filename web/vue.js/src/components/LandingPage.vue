<template>
  <div>
    <v-parallax :src="require('../assets/vbanner.jpg')">
      <v-layout class="center-text" align-center column justify-center>
        <h1 class="display-2 font-weight-thin mb-3">Simple Social Network</h1>
        <v-btn to="/signup" outline color="white">Sign up</v-btn>
      </v-layout>
    </v-parallax>

    <v-footer dark height="auto">
      <v-card flat tile class="primary white--text text-xs-center">
        <v-card-text>
          <v-btn v-for="icon in icons" :key="icon.icon" class="mx-3 white--text" :href="icon.link" target="_blank" icon>
            <v-icon size="24px">{{icon.icon}}</v-icon>
          </v-btn>
        </v-card-text>
        <v-card-text class="white--text pt-0">
          <LoremIpsum />
        </v-card-text>
        <v-divider></v-divider>
        <v-card-text class="white--text">
          &copy;2021 — <strong>Al Ol</strong>
        </v-card-text>
      </v-card>
    </v-footer>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

import LoremIpsum from '../components/LoremIpsum.vue';

@Component({
  components: {
    LoremIpsum,
  },
})
export default class LandingPage extends Vue {
  private icons = [
    { icon: 'fab fa-linkedin-in', link: 'https://www.linkedin.com/in/aovlllo/' },
    { icon: 'fab fa-twitter', link: 'https://twitter.com/aovlllo' },
    { icon: 'fab fa-github', link: 'https://github.com/aovlllo' },
  ];
}
</script>

<style scoped>
.center-text {
  text-align: center;
}

.center-image {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

.padding-top-bottom {
  padding-top: 50px;
  padding-bottom: 50px;
}
</style>
