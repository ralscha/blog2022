import {Component, inject} from '@angular/core';
import sqlite3InitModule from '@sqlite.org/sqlite-wasm';
import {HttpClient} from "@angular/common/http";
import {env, pipeline, TextGenerationPipeline} from "@huggingface/transformers";

type Country = {
  id: number,
  name: string,
  area: number,
  area_land: number,
  area_water: number,
  population: number,
  population_growth: number,
  birth_rate: number,
  death_rate: number,
  migration_rate: number,
  flag_description: string,
}


@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrl: './home.page.scss',
})
export class HomePage {
  readonly httpClient = inject(HttpClient);

  selectStatement = 'select * from countries;';
  countries: Country[] = [];
  searchTerm = '';
  db: any | undefined;
  generator: TextGenerationPipeline | undefined;
  working = false;

  readonly #prompt_template = `
    You are given a database schema and a question.
    Based on the schema, you need to generate a valid SQL SELECT query for sqlite that answers the question.

    Schema:
    CREATE TABLE countries (
        id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
        name TEXT NOT NULL,
        area       INTEGER,
        area_land  INTEGER,
        area_water INTEGER,
        population        INTEGER,
        population_growth REAL,
        birth_rate        REAL,
        death_rate        REAL,
        migration_rate    REAL,
        flag_description TEXT
      )

    Based on the above schema, generate a SQL SELECT query for the following question:

    Question: {question}

    Generate the SQL query based on the schema and the question. The query should always start with "SELECT * FROM countries"
    Return only the query. Do not include any comments or extra whitespace.
    `;

  constructor() {
    const httpClient = this.httpClient;

    const start = (sqlite3: any) => {
      httpClient.get('assets/countries.sqlite3', {responseType: 'arraybuffer'}).subscribe((data) => {
        const p = sqlite3.wasm.allocFromTypedArray(data);
        this.db = new sqlite3.oo1.DB();
        const deserialize_flags = sqlite3.capi.SQLITE_DESERIALIZE_FREEONCLOSE;
        const rc = sqlite3.capi.sqlite3_deserialize(this.db.pointer, 'main', p, data.byteLength, data.byteLength, deserialize_flags);
        this.db.checkRc(rc);

        this.countries = [];
        this.db.exec({
          sql: 'select * from countries;',
          rowMode: "object",
          callback: (row: any) => {
            this.countries.push(row);
          }
        });
      });
    };

    const initializeSQLite = async () => {
      try {
        const sqlite3 = await sqlite3InitModule({
          locateFile: (file) => `/sqlite-wasm/${file}`,
          print: console.log,
          printErr: console.error,
        });
        start(sqlite3);
      } catch (err) {
        console.error('Initialization error:', err);
      }
    };

    initializeSQLite();
    this.#initializeLLM();

  }

  async generateSQL(): Promise<void> {
    if (!this.generator || !this.db) {
      return;
    }
    this.countries = [];
    this.working = true;
    const userPrompt = this.#prompt_template.replace('{question}', this.searchTerm);
    this.selectStatement = '';

    const messages = [
      {role: "user", content: userPrompt},
    ];

    const output: any = await this.generator(messages, {max_new_tokens: 200});
    this.selectStatement = output[0].generated_text.at(-1).content;

    this.working = false;

    this.db.exec({
      sql: this.selectStatement,
      rowMode: "object",
      callback: (row: any) => {
        this.countries.push(row);
      }
    });
  }

  async isWebGPUAvailable(): Promise<boolean> {
    // @ts-expect-error
    if (!navigator.gpu) {
      return false;
    }
    try {
      // @ts-expect-error
      const adapter = await navigator.gpu.requestAdapter();
      return !!adapter;
    } catch {
      return false;
    }
  }

  async #initializeLLM() {
    env.localModelPath = '/assets';
    env.allowLocalModels = true;
    env.allowRemoteModels = false;
    env.backends.onnx.wasm!.wasmPaths = 'transformers-wasm/';

    this.generator = await pipeline('text-generation', 'onnx-community/Llama-3.2-1B-Instruct-q4f16', {
      device: 'webgpu',
      dtype: 'q4f16',
      local_files_only: true,
    });
  }
}
