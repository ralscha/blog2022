import {
  AutoProcessor,
  AutoTokenizer,
  MoonshineForConditionalGeneration,
  PreTrainedModel,
  PreTrainedTokenizer,
  Processor,
  TextStreamer,
} from '@huggingface/transformers';


class AutomaticSpeechRecognitionPipeline {
  static model_id: string | null = null;
  static tokenizer: Promise<PreTrainedTokenizer> | null = null;
  static processor: Promise<Processor> | null = null;
  static model: Promise<PreTrainedModel> | null = null;

  static async getInstance() {
    this.model_id = 'onnx-community/moonshine-tiny-ONNX';

    this.tokenizer = AutoTokenizer.from_pretrained(this.model_id, {});
    this.processor = AutoProcessor.from_pretrained(this.model_id, {});

    this.model = MoonshineForConditionalGeneration.from_pretrained(this.model_id, {
      dtype: {
        encoder_model: 'fp32',
        decoder_model_merged: 'q4',
      },
      device: 'webgpu',
    });
    return Promise.all([this.tokenizer, this.processor, this.model]);
  }
}

let processing = false;

// @ts-ignore
async function generate({audio, language}) {
  if (processing) {
    return;
  }
  processing = true;

  const [tokenizer, processor, model] = await AutomaticSpeechRecognitionPipeline.getInstance();

  const streamer = new TextStreamer(tokenizer, {
    skip_prompt: true,
    decode_kwargs: {
      skip_special_tokens: true,
    }
  });

  const inputs = await processor(audio);

  const outputs: any = await model.generate({
    ...inputs,
    max_new_tokens: 64,
    language,
    streamer,
  });

  const outputText = tokenizer.batch_decode(outputs, {skip_special_tokens: true});

  self.postMessage({
    status: 'complete',
    output: outputText,
  });
  processing = false;
}

async function load() {
  await AutomaticSpeechRecognitionPipeline.getInstance();
  self.postMessage({status: 'ready'});
}

self.addEventListener('message', async (e: MessageEvent) => {
  const {type, data} = e.data;

  switch (type) {
    case 'load':
      load();
      break;

    case 'generate':
      generate(data);
      break;
  }
});
