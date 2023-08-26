<template>
  <div>
    <div>
      <input v-model="searchTerm" placeholder="Search"/>
    </div>
    <table>
      <thead>
      <tr>
        <th>Name</th>
        <th>Age</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="row in filteredRows.slice((currentPage - 1) * pageSize, currentPage * pageSize)" :key="row.id">
        <td>{{ row.name }}</td>
        <td>{{ row.age }}</td>
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
    const tableData = ref({rows: []});
    const searchTerm = ref("");
    const currentPage = ref(1);
    const pageSize = ref(50);

    // Protobuf schema hardcoded as a string
    const protoSchema = `
      syntax = "proto3";
      message TableData {
        repeated Row rows = 1;
      }
      message Row {
        string name = 1;
        int32 age = 2;
        int64 id = 3;
      }
    `;

    // Parse the .proto schema
    const root = protobuf.parse(protoSchema).root;
    const TableData = root.lookupType("TableData");

    // WebSocket connection
    console.log("Connecting to WebSocket");
    const socket = new WebSocket("ws://localhost:8080/data");

    socket.onmessage = (event) => {
      const reader = new FileReader();
      reader.onload = () => {
        const buffer = new Uint8Array(reader.result);
        const value = TableData.decode(buffer);
        console.log("Decoded value:", value);
        tableData.value.rows = [...value.rows, ...tableData.value.rows];
      };
      reader.readAsArrayBuffer(event.data);
    };


    // Filtered rows based on search term
    const filteredRows = computed(() => {
      return tableData.value.rows.filter((row) =>
          row.name.toLowerCase().includes(searchTerm.value.toLowerCase())); // Note: Updated to "name"
    });

    // Pagination
    const pageCount = computed(() => {
      return Math.ceil(filteredRows.value.length / pageSize.value);
    });

    return {
      tableData,
      searchTerm,
      currentPage,
      pageSize,
      pageCount,
      filteredRows,
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
