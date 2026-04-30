import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { ParticleOptions } from '../types/pixi';

interface Particle {
  sprite: PIXI.Graphics;
  velocity: { x: number; y: number };
  life: number;
  maxLife: number;
}

export class MusicParticleSystem {
  private container: PIXI.Container;
  private particles: Particle[] = [];
  private particlePool: PIXI.Graphics[] = [];
  private maxParticles: number = 1000;

  constructor(container: PIXI.Container, maxParticles: number = 1000) {
    this.container = container;
    this.maxParticles = maxParticles;
    
    // Pre-create particle pool for performance
    this.initializePool();
  }

  /**
   * Pre-create particles for object pooling
   */
  private initializePool(): void {
    for (let i = 0; i < this.maxParticles; i++) {
      const particle = new PIXI.Graphics();
      particle.visible = false;
      this.container.addChild(particle);
      this.particlePool.push(particle);
    }
  }

  /**
   * Get particle from pool
   */
  private getParticle(): PIXI.Graphics | null {
    return this.particlePool.find(p => !p.visible) || null;
  }

  /**
   * Create particle burst (for note events)
   */
  createBurst(
    x: number,
    y: number,
    options: Partial<ParticleOptions> = {}
  ): void {
    const opts = {
      color: options.color || 0x3B82F6,
      size: options.size || 3,
      lifetime: options.lifetime || 1,
      velocity: options.velocity || { x: 0, y: -50 },
      gravity: options.gravity || 100
    };

    // Calculate particle count based on intensity
    const count = Math.floor(10 + Math.random() * 20);

    for (let i = 0; i < count; i++) {
      const particle = this.getParticle();
      if (!particle) continue;

      // Random angle
      const angle = Math.random() * Math.PI * 2;
      const speed = 50 + Math.random() * 100;
      
      const velocity = {
        x: Math.cos(angle) * speed,
        y: Math.sin(angle) * speed
      };

      // Draw particle
      particle.clear();
      particle.beginFill(opts.color);
      particle.drawCircle(0, 0, opts.size);
      particle.endFill();
      
      particle.x = x;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      // Add to active particles
      this.particles.push({
        sprite: particle,
        velocity,
        life: opts.lifetime,
        maxLife: opts.lifetime
      });
    }
  }

  /**
   * Create continuous stream (for sustained notes)
   */
  createStream(
    x: number,
    y: number,
    color: number,
    intensity: number = 1
  ): void {
    const count = Math.floor(intensity * 5);
    
    for (let i = 0; i < count; i++) {
      const particle = this.getParticle();
      if (!particle) continue;

      // Upward velocity with slight randomness
      const velocity = {
        x: (Math.random() - 0.5) * 30,
        y: -50 - Math.random() * 50
      };

      // Draw particle
      particle.clear();
      particle.beginFill(color, 0.8);
      particle.drawCircle(0, 0, 2 + Math.random() * 2);
      particle.endFill();
      
      particle.x = x + (Math.random() - 0.5) * 20;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      this.particles.push({
        sprite: particle,
        velocity,
        life: 0.5 + Math.random() * 0.5,
        maxLife: 1
      });
    }
  }

  /**
   * Create explosion (for beat hits)
   */
  createExplosion(
    x: number,
    y: number,
    color: number,
    size: number = 50
  ): void {
    const count = Math.floor(30 + size);

    for (let i = 0; i < count; i++) {
      const particle = this.getParticle();
      if (!particle) continue;

      const angle = (i / count) * Math.PI * 2;
      const speed = size * 2 + Math.random() * size;
      
      const velocity = {
        x: Math.cos(angle) * speed,
        y: Math.sin(angle) * speed
      };

      // Draw particle with glow
      particle.clear();
      particle.beginFill(color, 0.9);
      particle.drawCircle(0, 0, 2 + Math.random() * 3);
      particle.endFill();
      
      // Add glow filter
      particle.filters = [new PIXI.BlurFilter({ strength: 2 })];
      
      particle.x = x;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      this.particles.push({
        sprite: particle,
        velocity,
        life: 0.8 + Math.random() * 0.4,
        maxLife: 1.2
      });
    }
  }

  /**
   * Create wave pattern (for harmonics)
   */
  createWave(
    startX: number,
    startY: number,
    endX: number,
    endY: number,
    color: number,
    amplitude: number = 50
  ): void {
    const distance = Math.sqrt(
      Math.pow(endX - startX, 2) + 
      Math.pow(endY - startY, 2)
    );
    const steps = Math.floor(distance / 10);

    for (let i = 0; i < steps; i++) {
      const t = i / steps;
      const particle = this.getParticle();
      if (!particle) continue;

      // Sine wave offset
      const waveOffset = Math.sin(t * Math.PI * 4) * amplitude;
      
      const x = startX + (endX - startX) * t;
      const y = startY + (endY - startY) * t + waveOffset;

      // Draw particle
      particle.clear();
      particle.beginFill(color, 0.6);
      particle.drawCircle(0, 0, 2);
      particle.endFill();
      
      particle.x = x;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      this.particles.push({
        sprite: particle,
        velocity: { x: 0, y: 20 },
        life: 1 + t * 0.5,
        maxLife: 1.5
      });
    }
  }

  /**
   * Update particles (call in animation loop)
   */
  update(deltaTime: number): void {
    for (let i = this.particles.length - 1; i >= 0; i--) {
      const particle = this.particles[i];
      
      // Update life
      particle.life -= deltaTime;
      
      if (particle.life <= 0) {
        // Return to pool
        particle.sprite.visible = false;
        particle.sprite.filters = [];
        this.particles.splice(i, 1);
        continue;
      }

      // Update position
      particle.sprite.x += particle.velocity.x * deltaTime;
      particle.sprite.y += particle.velocity.y * deltaTime;

      // Apply gravity
      particle.velocity.y += 100 * deltaTime;

      // Fade out based on remaining life
      particle.sprite.alpha = particle.life / particle.maxLife;
    }
  }

  /**
   * Clear all particles
   */
  clear(): void {
    this.particles.forEach(p => {
      p.sprite.visible = false;
      p.sprite.filters = [];
    });
    this.particles = [];
  }

  /**
   * Cleanup
   */
  destroy(): void {
    this.clear();
    this.particlePool.forEach(p => p.destroy());
    this.particlePool = [];
  }
}

