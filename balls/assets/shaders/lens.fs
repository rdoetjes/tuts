#version 330

// Input vertex attributes
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// Lens uniforms
uniform vec2 lensPos;
uniform float lensRadius;
uniform float magnification;
uniform vec2 renderSize;

void main()
{
    // Screen coordinates
    vec2 screenPos = fragTexCoord * renderSize;

    // Distance from pixel to lens center
    float dist = distance(screenPos, lensPos);

    if (dist < lensRadius)
    {
        // Magnification logic: Scale coordinates relative to the lens center
        vec2 offset = screenPos - lensPos;

        // Curved lens logic (Fish-eye effect)
        float falloff = 1.0 - (dist / lensRadius); // 1.0 at center, 0.0 at edge
        float dynamicMag = 1.0 + (magnification - 1.0) * falloff;
        vec2 magCoord = lensPos + offset / dynamicMag;

        // Normalize back to 0.0 - 1.0 for texture sampling
        vec2 uv = magCoord / renderSize;

        vec4 texel = texture(texture0, uv);

        // Simple border/rim effect
        float edge = smoothstep(lensRadius - 2.0, lensRadius, dist);
        vec4 borderColor = vec4(1.0, 1.0, 1.0, 0.3); // Subtle white rim

        finalColor = mix(texel, borderColor*3, edge);
    }
    else
    {
        finalColor = texture(texture0, fragTexCoord);
    }
}
