#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

uniform vec2 resolution;

// CRT effect parameters
const float warpX = 0.03; // Screen curvature X
const float warpY = 0.04; // Screen curvature Y
const float scanlineAlpha = 0.12;
const float vignetteIntensity = 0.15;

vec2 curve(vec2 uv) {
    uv = uv * 2.0 - 1.0;
    uv.x *= 1.0 + (uv.y * uv.y) * warpX;
    uv.y *= 1.0 + (uv.x * uv.x) * warpY;
    uv = uv * 0.5 + 0.5;
    return uv;
}

void main() {
    vec2 uv = curve(fragTexCoord);

    // Check if we're out of bounds after curvature
    if (uv.x < 0.0 || uv.x > 1.0 || uv.y < 0.0 || uv.y > 1.0) {
        finalColor = vec4(0.0, 0.0, 0.0, 1.0);
        return;
    }

    // Sample texture with slight chromatic aberration
    float distortion = 0.002;
    vec4 color;
    color.r = texture(texture0, vec2(uv.x + distortion, uv.y)).r;
    color.g = texture(texture0, uv).g;
    color.b = texture(texture0, vec2(uv.x - distortion, uv.y)).b;
    color.a = 1.0;

    // Scanlines
    float scanline = sin(uv.y * resolution.y * 1.5) * scanlineAlpha;
    color.rgb -= scanline;

    // Vignette
    float vignette = uv.x * uv.y * (1.0 - uv.x) * (1.0 - uv.y);
    vignette = pow(vignette * 15.0, vignetteIntensity);
    color.rgb *= vignette;

    finalColor = color * colDiffuse;
}
