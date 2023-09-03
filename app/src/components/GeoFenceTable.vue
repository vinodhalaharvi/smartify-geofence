<template>
  <div>
    <div>
      <input v-model="searchTerm" placeholder="Search by Polygon ID"/>
    </div>
    <table>
      <thead>
      <tr>
        <th>Latitude</th>
        <th>Longitude</th>
        <th>Polygon ID</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="location in filteredLocations.slice((currentPage - 1) * pageSize, currentPage * pageSize)"
          :key="location.id">
        <td>{{ location.latitude }}</td>
        <td>{{ location.longitude }}</td>
        <td>{{ location.polygonId }}</td>
      </tr>
      </tbody>
    </table>
    <div>
      <button @click="currentPage--" :disabled="currentPage <= 1">Previous</button>
      <button @click="currentPage++" :disabled="currentPage >= pageCount">Next</button>
    </div>
  </div>
</template>

<script>
import protobuf from "protobufjs";
import {computed, ref} from "vue";

export default {
  setup() {
    const locations = ref([]);
    const searchTerm = ref("");
    const currentPage = ref(1);
    const pageSize = ref(50);

    const protoSchema = `
      syntax = "proto3";
      message FencedLocation {
        double latitude = 1;
        double longitude = 2;
        string polygonId = 3;
      }
    `;

    const root = protobuf.parse(protoSchema).root;
    const FencedLocation = root.lookupType("FencedLocation");

    console.log("Connecting to WebSocket");
    const socket = new WebSocket("ws://localhost:8080/data");

    socket.onmessage = (event) => {
      const reader = new FileReader();
      reader.onload = () => {
        const buffer = new Uint8Array(reader.result);
        const value = FencedLocation.decode(buffer);
        console.log("Decoded value:", value);
        locations.value = [value, ...locations.value];
      };
      reader.readAsArrayBuffer(event.data);
    };

    const filteredLocations = computed(() => {
      return locations.value.filter((location) =>
          location.polygonId.toLowerCase().includes(searchTerm.value.toLowerCase())
      );
    });

    const pageCount = computed(() => {
      return Math.ceil(filteredLocations.value.length / pageSize.value);
    });

    return {
      locations,
      searchTerm,
      currentPage,
      pageSize,
      pageCount,
      filteredLocations,
    };
  },
};
</script>

<style scoped>
div {
  margin: 20px;
}

input {
  margin-bottom: 10px;
  padding: 5px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 15px;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}

tr:nth-child(even) {
  background-color: #f2f2f2;
}

tr:hover {
  background-color: #f5f5f5;
}

button {
  padding: 10px 20px;
  margin: 5px;
  cursor: pointer;
  background-color: #ec6868;
  color: white;
  border: none;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
</style>
