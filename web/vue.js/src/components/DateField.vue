<template>
  <v-text-field v-model="birth" :rules="dateRules" prepend-icon="calendar_today" name="birth" label="Birth date" type="date" required></v-text-field>
</template>

<script lang="ts">
import { Component, Emit, Prop, Vue } from 'vue-property-decorator';

@Component
export default class DateField extends Vue {
  @Prop(String) public value!: string;

  private dateRules = [
    (v: string) => !!v || 'Birth date is required',
  ];

  private get birth(): string {
    return this.value;
  }

  private set birth(newVal: string) {
    this.changeValue(newVal);
  }

  @Emit('input')
  public changeValue(val: string) { /* */ }
}
</script>
