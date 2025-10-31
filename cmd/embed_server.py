# cmd/embed_server.py
from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer

app = Flask(__name__)
model = SentenceTransformer('BAAI/bge-m3')

@app.route('/embed', methods=['POST'])
def embed():
    data = request.json
    texts = data.get('texts', [])
    if not isinstance(texts, list):
        return jsonify({"error": "texts must be a list"}), 400
    
    # ✅ 关键：启用归一化（bge-m3 官方推荐）
    embeddings = model.encode(
        texts,
        convert_to_numpy=True,
        normalize_embeddings=True  # ← 添加这一行
    ).tolist()
    
    return jsonify({"embeddings": embeddings})

if __name__ == '__main__':
    app.run(host='127.0.0.1', port=8081)