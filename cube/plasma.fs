#version 330

// Using gl_FragCoord instead of fragTexCoord to ensure we see the "whole area"
// correctly across the entire 800x600 screen.
uniform float uTime;
uniform vec2 uResolution;

out vec4 finalColor;

void main()
{
    // Normalized pixel coordinates (from 0 to 1)
    vec2 uv = gl_FragCoord.xy / uResolution.xy;
    
    // Zoom in more by reducing the coordinate multipliers
    float fx = uv.x * 240.0;
    float fy = uv.y * 180.0;
    
    // Slowed down by another 200% (3x slower than before)
    float t = uTime * 0.4;

    // We add many more sine waves with different frequencies and phases
    // to create a dense, "demo-scene" interference pattern.
    float v = 0.0;
    
    // 1. Denser base waves
    v += sin(fx / 10.0 + t);
    v += sin(fy / 15.0 + t * 0.8);
    v += sin((fx + fy + t * 5.0) / 20.0);
    
    // 2. Multiple radial interference points (The Swirls)
    // Moving the centers around makes them "swirl" into each other.
    v += sin(sqrt(pow(fx - 100.0 + 50.0 * sin(t), 2.0) + pow(fy - 75.0 + 50.0 * cos(t), 2.0)) / 8.0);
    v += sin(sqrt(pow(fx - 300.0 + 80.0 * cos(t * 0.5), 2.0) + pow(fy - 150.0 + 60.0 * sin(t * 0.7), 2.0)) / 10.0);
    v += sin(sqrt(pow(fx - 200.0, 2.0) + pow(fy - 250.0, 2.0)) / 12.0 + t);
    
    // 3. High-frequency diagonal interference
    v += sin((fx * 0.5 - fy * 0.3 + t * 2.0) / 5.0);

    // Color mapping with high frequency cycles
    float pi = 3.14159265;
    
    // Multiply v by pi * 2.0 (or more) to get many distinct color bands
    float r = 0.5 + 0.5 * sin(v * pi);
    float g = 0.5 + 0.5 * sin(v * pi + 2.0);
    float b = 0.5 + 0.5 * sin(v * pi + 4.0);

    finalColor = vec4(r, g, b, 1.0);
}
