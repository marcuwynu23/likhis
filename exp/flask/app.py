from flask import Flask, request, jsonify

app = Flask(__name__)

# Root route
@app.route('/')
def index():
    return jsonify({'message': 'Welcome to the API'})

# Health check
@app.route('/health')
def health():
    return jsonify({'status': 'ok'})

# Example route with query parameters
@app.route('/search')
def search():
    q = request.args.get('q')
    page = request.args.get('page')
    return jsonify({'query': q, 'page': page})

# Example POST route
@app.route('/auth/login', methods=['POST'])
def login():
    email = request.form.get('email')
    password = request.form.get('password')
    return jsonify({'message': 'Login successful'})

# Example route with path parameter
@app.route('/users/<int:user_id>')
def get_user(user_id):
    return jsonify({'id': user_id})

# Example route with string parameter
@app.route('/products/<name>')
def get_product(name):
    return jsonify({'name': name})

# Example PUT route
@app.route('/settings/<int:user_id>', methods=['PUT'])
def update_settings(user_id):
    theme = request.form.get('theme')
    return jsonify({'userId': user_id, 'theme': theme})

# Example DELETE route
@app.route('/sessions/<session_id>', methods=['DELETE'])
def delete_session(session_id):
    return jsonify({'message': 'Session deleted', 'sessionId': session_id})

if __name__ == '__main__':
    app.run(debug=True)

