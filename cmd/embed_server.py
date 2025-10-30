# embed_server.py
from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer

app = Flask(__name__)

# 加载 all-MiniLM-L6-v2（首次运行会自动下载，约 80MB）
model = SentenceTransformer('all-MiniLM-L6-v2')

@app.route('/embed', methods=['POST'])
def embed():
    data = request.json
    texts = data.get('texts', [])
    if not isinstance(texts, list):
        return jsonify({"error": "texts must be a list of strings"}), 400
    embeddings = model.encode(texts, convert_to_numpy=True).tolist()
    return jsonify({"embeddings": embeddings})

if __name__ == '__main__':
    app.run(host='127.0.0.1', port=8081)